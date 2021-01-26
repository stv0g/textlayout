package harfbuzz

import (
	"fmt"
	"sync"
)

// ported from harfbuzz/src/hb-shape.cc, harfbuzz/src/hb-shape-plan.cc
// Copyright © 2009  Red Hat, Inc.
// Copyright © 2012  Google, Inc.
// Red Hat Author(s): Behdad Esfahbod
// Google Author(s): Behdad Esfahbod

/**
 * SECTION:hb-shape
 * @title: hb-shape
 * @short_description: Conversion of text strings into positioned glyphs
 * @include: hb.h
 *
 * Shaping is the central operation of HarfBuzz. Shaping operates on buffers,
 * which are sequences of Unicode characters that use the same font and have
 * the same text direction, script, and language. After shaping the buffer
 * contains the output glyphs and their positions.
 **/

//  static const char *nil_shaper_list[] = {nullptr};

func create() []string {
	var out []string
	shapers := _hb_shapers_get()
	for _, s := range shapers {
		out = append(out, s.name)
	}
	return out
}

// See hb_shape() for details. If `shaperList` is not empty, the specified
// shapers will be used in the given order, otherwise the default shapers list
// will be used.
//
// Return value: false if all shapers failed, true otherwise
func (buffer *hb_buffer_t) hb_shape_full(font *hb_font_t, features []hb_feature_t, shaperList []string) bool {
	shape_plan := hb_shape_plan_create_cached2(font.face, &buffer.props,
		features, font.coords, shaperList)
	res := hb_shape_plan_execute(shape_plan, font, buffer, features)
	return res
}

// Shapes `buffer` using `font`, turning its Unicode characters content to
// positioned glyphs. If `features` is not empty, it will be used to control the
// features applied during shaping. If two features have the same tag but
// overlapping ranges the value of the feature with the higher index takes
// precedence.
func (buffer *hb_buffer_t) hb_shape(font *hb_font_t, features []hb_feature_t) {
	buffer.hb_shape_full(font, features, nil)
}

// use interface since equality check is needed
type hb_shape_func_t = func(shape_plan *hb_shape_plan_t,
	font *hb_font_t, buffer *hb_buffer_t, features []hb_feature_t) bool

type hb_ot_shape_plan_key_t = [2]int

func hb_ot_shape_plan_key_t_init(face hb_face_t, coords []float32) hb_ot_shape_plan_key_t {
	gsub := face.hb_ot_layout_table_find_feature_variations(HB_OT_TAG_GSUB, coords)
	gpos := face.hb_ot_layout_table_find_feature_variations(HB_OT_TAG_GPOS, coords)
	return [2]int{gsub, gpos}
}

type hb_shape_plan_key_t struct {
	props hb_segment_properties_t

	userFeatures []hb_feature_t // len num_user_features

	ot hb_ot_shape_plan_key_t

	shaper_func hb_shape_func_t
	shaper_name string
}

/*
 * hb_shape_plan_key_t
 */

func (plan *hb_shape_plan_key_t) init(copy bool,
	face hb_face_t, props *hb_segment_properties_t,
	userFeatures []hb_feature_t, coords []float32, shaperList []string) {

	plan.props = *props
	if !copy {
		plan.userFeatures = userFeatures
	} else {
		plan.userFeatures = append([]hb_feature_t(nil), userFeatures...)
		/* Make start/end uniform to easier catch bugs. */
		for i := range plan.userFeatures {
			if plan.userFeatures[i].start != HB_FEATURE_GLOBAL_START {
				plan.userFeatures[i].start = 1
			}
			if plan.userFeatures[i].end != HB_FEATURE_GLOBAL_END {
				plan.userFeatures[i].end = 2
			}
		}
	}
	plan.ot = hb_ot_shape_plan_key_t_init(face, coords)
}

func (plan hb_shape_plan_key_t) user_features_match(other hb_shape_plan_key_t) bool {
	if len(plan.userFeatures) != len(other.userFeatures) {
		return false
	}
	for i, feat := range plan.userFeatures {
		if feat.tag != other.userFeatures[i].tag || feat.value != other.userFeatures[i].value ||
			(feat.start == HB_FEATURE_GLOBAL_START && feat.end == HB_FEATURE_GLOBAL_END) !=
				(other.userFeatures[i].start == HB_FEATURE_GLOBAL_START && other.userFeatures[i].end == HB_FEATURE_GLOBAL_END) {
			return false
		}
	}
	return true
}

func (plan hb_shape_plan_key_t) equal(other hb_shape_plan_key_t) bool {
	return plan.props == other.props &&
		plan.user_features_match(other) &&
		plan.ot == other.ot && plan.shaper_func == other.shaper_func
}

// Shape plans are an internal mechanism. Each plan contains state
// describing how HarfBuzz will shape a particular text segment, based on
// the combination of segment properties and the capabilities in the
// font face in use.
//
// Shape plans are not used for shaping directly, but can be queried to
// access certain information about how shaping will perform, given a set
// of specific input parameters (script, language, direction, features,
// etc.).
//
// Most client programs will not need to deal with shape plans directly.
type hb_shape_plan_t struct {
	face_unsafe hb_face_t
	key         hb_shape_plan_key_t
	ot          hb_ot_shape_plan_t
}

/**
 * hb_shape_plan_create: (Xconstructor)
 * @face: #hb_face_t to use
 * @props: The #hb_segment_properties_t of the segment
 * @userFeatures: (array length=num_user_features): The list of user-selected features
 * @num_user_features: The number of user-selected features
 * @shaperList: (array zero-terminated=1): List of shapers to try
 *
 * Constructs a shaping plan for a combination of @face, @userFeatures, @props,
 * and @shaperList.
 *
 * Return value: (transfer full): The shaping plan
 *
 * Since: 0.9.7
 **/
func hb_shape_plan_create(face hb_face_t, props *hb_segment_properties_t,
	userFeatures []hb_feature_t, shaperList []string) *hb_shape_plan_t {
	return hb_shape_plan_create2(face, props, userFeatures, nil, shaperList)
}

/**
 * hb_shape_plan_create2: (Xconstructor)
 * @face: #hb_face_t to use
 * @props: The #hb_segment_properties_t of the segment
 * @userFeatures: (array length=num_user_features): The list of user-selected features
 * @num_user_features: The number of user-selected features
 * @coords: (array length=num_coords): The list of variation-space coordinates
 * @num_coords: The number of variation-space coordinates
 * @shaperList: (array zero-terminated=1): List of shapers to try
 *
 * The variable-font version of #hb_shape_plan_create.
 * Constructs a shaping plan for a combination of @face, @userFeatures, @props,
 * and @shaperList, plus the variation-space coordinates @coords.
 *
 * Return value: (transfer full): The shaping plan
 *
 * Since: 1.4.0
 **/

func hb_shape_plan_create2(face hb_face_t, props *hb_segment_properties_t,
	userFeatures []hb_feature_t, coords []float32, shaperList []string) *hb_shape_plan_t {

	if debugMode {
		fmt.Printf("shape plan: face:%p num_features:%d num_coords=%d shaperList:%s", face, len(userFeatures), len(coords), shaperList)
	}

	var shape_plan hb_shape_plan_t

	shape_plan.face_unsafe = face

	shape_plan.key.init(true, face, props, userFeatures, coords, shaperList)
	shape_plan.ot.init0(face, &shape_plan.key)

	return &shape_plan

}

/*
* Choose shaper.
 */

//  #define HB_SHAPER_PLAN(shaper) \
// 	 HB_STMT_START { \
// 	   if (face.data.shaper) \
// 	   { \
// 		 this.shaper_func = _hb_##shaper##_shape; \
// 		 this.shaper_name = #shaper; \
// 		 return true; \
// 	   } \
// 	 } HB_STMT_END

//    if (unlikely (shaperList))
//    {
// 	 for (; *shaperList; shaperList++)
// 	   if (false)
// 	 ;
//  #define HB_SHAPER_IMPLEMENT(shaper) \
// 	   else if (0 == strcmp (*shaperList, #shaper)) \
// 	 HB_SHAPER_PLAN (shaper);
//  #include "hb-shaper-list.hh"
//  #undef HB_SHAPER_IMPLEMENT
//    }
//    else
//    {
// 	 const hb_shaper_entry_t *shapers = _hb_shapers_get ();
// 	 for (unsigned int i = 0; i < HB_SHAPERS_COUNT; i++)
// 	   if (false)
// 	 ;
//  #define HB_SHAPER_IMPLEMENT(shaper) \
// 	   else if (shapers[i].func == _hb_##shaper##_shape) \
// 	 HB_SHAPER_PLAN (shaper);
//  #include "hb-shaper-list.hh"
//  #undef HB_SHAPER_IMPLEMENT
//    }
//  #undef HB_SHAPER_PLAN

//  bail:
//    ::free (features);
//    return false;
//  }

/**
 * hb_shape_plan_get_empty:
 *
 * Fetches the singleton empty shaping plan.
 *
 * Return value: (transfer full): The empty shaping plan
 *
 * Since: 0.9.7
 **/
//  hb_shape_plan_t *
//  hb_shape_plan_get_empty ()
//  {
//    return const_cast<hb_shape_plan_t *> (&Null (hb_shape_plan_t));
//  }

//  static bool
//  _hb_shape_plan_execute_internal (shape_plan *hb_shape_plan_t,
// 				  hb_font_t          *font,
// 				  hb_buffer_t        *buffer,
// 				  const hb_feature_t *features,
// 				  unsigned int        num_features)
//  {
//    DEBUG_MSG_FUNC (SHAPE_PLAN, shape_plan,
// 		   "num_features=%d shaper_func=%p, shaper_name=%s",
// 		   num_features,
// 		   shape_plan.key.shaper_func,
// 		   shape_plan.key.shaper_name);

//    if (unlikely (!buffer.len))
// 	 return true;

//    assert (!hb_object_is_immutable (buffer));

//    buffer.assert_unicode ();

//    if (unlikely (hb_object_is_inert (shape_plan)))
// 	 return false;

//    assert (shape_plan.face_unsafe == font.face);
//    assert (hb_segment_properties_equal (&shape_plan.key.props, &buffer.props));

//  #define HB_SHAPER_EXECUTE(shaper) \
// 	 HB_STMT_START { \
// 	   return font.data.shaper && \
// 		  _hb_##shaper##_shape (shape_plan, font, buffer, features, num_features); \
// 	 } HB_STMT_END

//    if (false)
// 	 ;
//  #define HB_SHAPER_IMPLEMENT(shaper) \
//    else if (shape_plan.key.shaper_func == _hb_##shaper##_shape) \
// 	 HB_SHAPER_EXECUTE (shaper);
//  #include "hb-shaper-list.hh"
//  #undef HB_SHAPER_IMPLEMENT

//  #undef HB_SHAPER_EXECUTE

//    return false;
//  }

/**
 * hb_shape_plan_execute:
 * @shape_plan: A shaping plan
 * @font: The #hb_font_t to use
 * @buffer: The #hb_buffer_t to work upon
 * @features: (array length=num_features): Features to enable
 * @num_features: The number of features to enable
 *
 * Executes the given shaping plan on the specified buffer, using
 * the given @font and @features.
 *
 * Return value: %true if success, %false otherwise.
 *
 * Since: 0.9.7
 **/
func (shape_plan *hb_shape_plan_t) hb_shape_plan_execute(
	font *hb_font_t, buffer *hb_buffer_t,
	features []hb_feature_t) bool {
	ret := _hb_shape_plan_execute_internal(shape_plan, font, buffer, features)

	if ret && buffer.content_type == HB_BUFFER_CONTENT_TYPE_UNICODE {
		buffer.content_type = HB_BUFFER_CONTENT_TYPE_GLYPHS
	}

	return ret
}

/*
 * Caching
 */

/**
 * hb_shape_plan_create_cached:
 * @face: #hb_face_t to use
 * @props: The #hb_segment_properties_t of the segment
 * @userFeatures: (array length=num_user_features): The list of user-selected features
 * @num_user_features: The number of user-selected features
 * @shaperList: (array zero-terminated=1): List of shapers to try
 *
 * Creates a cached shaping plan suitable for reuse, for a combination
 * of @face, @userFeatures, @props, and @shaperList.
 *
 * Return value: (transfer full): The shaping plan
 *
 * Since: 0.9.7
 **/
func hb_shape_plan_create_cached(face hb_face_t,
	props *hb_segment_properties_t,
	userFeatures []hb_feature_t, shaperList []string) *hb_shape_plan_t {
	return hb_shape_plan_create_cached2(face, props, userFeatures, nil, shaperList)
}

var (
	planCache     map[hb_face_t][]*hb_shape_plan_t
	planCacheLock sync.Mutex
)

/**
 * hb_shape_plan_create_cached2:
 * @face: #hb_face_t to use
 * @props: The #hb_segment_properties_t of the segment
 * @userFeatures: (array length=num_user_features): The list of user-selected features
 * @num_user_features: The number of user-selected features
 * @coords: (array length=num_coords): The list of variation-space coordinates
 * @num_coords: The number of variation-space coordinates
 * @shaperList: (array zero-terminated=1): List of shapers to try
 *
 * The variable-font version of #hb_shape_plan_create_cached.
 * Creates a cached shaping plan suitable for reuse, for a combination
 * of @face, @userFeatures, @props, and @shaperList, plus the
 * variation-space coordinates @coords.
 *
 * Return value: (transfer full): The shaping plan
 *
 * Since: 1.4.0
 **/
func hb_shape_plan_create_cached2(face hb_face_t,
	props *hb_segment_properties_t,
	userFeatures []hb_feature_t, coords []float32, shaperList []string) *hb_shape_plan_t {
	if debugMode {
		fmt.Printf("shape plan: face:%p num_features:%d shaperList:%s", face, len(userFeatures), shaperList)
	}

	var key hb_shape_plan_key_t
	key.init(false, face, props, userFeatures, coords, shaperList)

	planCacheLock.Lock()
	defer planCacheLock.Unlock()

	plans := planCache[face]

	for _, plan := range plans {
		if plan.key.equal(key) {
			if debugMode {
				fmt.Println(plan, "fulfilled from cache")
			}
			return plan
		}
	}
	plan := hb_shape_plan_create2(face, props, userFeatures, coords, shaperList)

	plans = append(plans, plan)
	planCache[face] = plans
	if debugMode {
		fmt.Println(plan, "inserted into cache")
	}

	return plan
}
