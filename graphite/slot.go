package graphite

const (
	DELETED uint8 = 1 << iota
	INSERTED
	COPIED
	POSITIONED
	ATTACHED
)

// Slot represents one glyph in a line of text.
type Slot struct {
	prev, Next *Slot        // linked list of slots
	parent     *Slot        // index to parent we are attached to
	child      *Slot        // index to first child slot that attaches to us
	sibling    *Slot        // index to next child that attaches to our parent
	justs      *slotJustify // pointer to justification parameters

	userAttrs []int16 // with length silf.NumUserDefn

	original    int // charinfo that originated this slot (e.g. for feature values)
	Before      int // charinfo index of before association
	After       int // charinfo index of after association
	index       int // slot index given to this slot during finalising
	glyphID     GID
	realGlyphID GID
	Position    Position // absolute position of glyph
	shift       Position // .shift slot attribute
	advance     Position // .advance slot attribute
	attach      Position // attachment point on us
	with        Position // attachment point position on parent
	just        float32  // Justification inserted space
	flags       uint8    // holds bit flags
	attLevel    uint8    // attachment level
	bidiCls     int8     // bidirectional class
	bidiLevel   uint8    // bidirectional level
}

// returns true if the slot has no parent
func (sl *Slot) isBase() bool {
	return sl.parent == nil
}

// move up the tree and return the highest non nil slot
func (is *Slot) findRoot() *Slot {
	for ; is.parent != nil; is = is.parent {
	}
	return is
}

// return true if the slot has `base` in its ancesters
func (is *Slot) isChildOf(base *Slot) bool {
	for p := is.parent; p != nil; p = p.parent {
		if p == base {
			return true
		}
	}
	return false
}

func (sl *Slot) isDeleted() bool {
	return sl.flags&DELETED != 0
}

func (sl *Slot) markDeleted(state bool) {
	if state {
		sl.flags |= DELETED
	} else {
		sl.flags &= ^DELETED
	}
}

func (sl *Slot) isCopied() bool {
	return sl.flags&COPIED != 0
}

func (sl *Slot) markCopied(state bool) {
	if state {
		sl.flags |= COPIED
	} else {
		sl.flags &= ^COPIED
	}
}

func (sl *Slot) isInsertBefore() bool {
	return sl.flags&INSERTED != 0
}

func (sl *Slot) markInsertBefore(state bool) {
	if state {
		sl.flags |= INSERTED
	} else {
		sl.flags &= ^INSERTED
	}
}

func (sl *Slot) setGlyph(seg *Segment, glyphID GID) {
	sl.glyphID = glyphID
	sl.bidiCls = -1
	theGlyph := seg.face.getGlyph(glyphID)
	if theGlyph == nil {
		sl.realGlyphID = 0
		sl.advance = Position{}
		return

	}
	sl.realGlyphID = GID(theGlyph.attrs.get(uint16(seg.silf.attrPseudo)))
	if int(sl.realGlyphID) > len(seg.face.glyphs) {
		sl.realGlyphID = 0
	}
	aGlyph := theGlyph
	if sl.realGlyphID != 0 {
		aGlyph = seg.face.getGlyph(sl.realGlyphID)
		if aGlyph == nil {
			aGlyph = theGlyph
		}
	}
	sl.advance = Position{x: float32(aGlyph.advance.x), y: 0.}
	if seg.silf.attrSkipPasses != 0 {
		seg.mergePassBits(uint32(theGlyph.attrs.get(uint16(seg.silf.attrSkipPasses))))
		if len(seg.silf.passes) > 16 {
			seg.mergePassBits(uint32(theGlyph.attrs.get(uint16(seg.silf.attrSkipPasses)+1)) << 16)
		}
	}
}

func (sl *Slot) removeChild(ap *Slot) bool {
	if sl == ap || sl.child == nil || ap == nil {
		return false
	} else if ap == sl.child {
		nSibling := sl.child.sibling
		sl.child.sibling = nil
		sl.child = nSibling
		return true
	}
	for p := sl.child; p != nil; p = p.sibling {
		if p.sibling != nil && p.sibling == ap {
			p.sibling = p.sibling.sibling
			ap.sibling = nil
			return true
		}
	}
	return false
}

func (sl *Slot) setSibling(ap *Slot) bool {
	if sl == ap {
		return false
	} else if ap == sl.sibling {
		return true
	} else if sl.sibling == nil || ap == nil {
		sl.sibling = ap
	} else {
		return sl.sibling.setSibling(ap)
	}
	return true
}

func (sl *Slot) setChild(ap *Slot) bool {
	if sl == ap {
		return false
	} else if ap == sl.child {
		return true
	} else if sl.child == nil {
		sl.child = ap
	} else {
		return sl.child.setSibling(ap)
	}
	return true
}

func (sl *Slot) getJustify(seg *Segment, level uint8, subindex int) int16 {
	if level != 0 && int(level) >= len(seg.silf.justificationLevels) {
		return 0
	}

	if sl.justs != nil {
		return sl.justs.values[level][subindex]
	}

	if int(level) >= len(seg.silf.justificationLevels) {
		return 0
	}
	jAttrs := seg.silf.justificationLevels[level]

	switch subindex {
	case 0:
		return seg.face.getGlyphAttr(sl.glyphID, uint16(jAttrs.attrStretch))
	case 1:
		return seg.face.getGlyphAttr(sl.glyphID, uint16(jAttrs.attrShrink))
	case 2:
		return seg.face.getGlyphAttr(sl.glyphID, uint16(jAttrs.attrStep))
	case 3:
		return seg.face.getGlyphAttr(sl.glyphID, uint16(jAttrs.attrWeight))
	case 4:
		return 0 // not been set yet, so clearly 0
	}
	return 0
}

// #define SLOTGETCOLATTR(x) { SlotCollision *c = seg.collisionInfo(this); return c ? int(c. x) : 0; }

func (sl *Slot) getAttr(seg *Segment, ind attrCode, subindex int) int32 {
	if ind >= gr_slatJStretch && ind < gr_slatJStretch+20 && ind != gr_slatJWidth {
		indx := int(ind - gr_slatJStretch)
		return int32(sl.getJustify(seg, uint8(indx/NUMJUSTPARAMS), indx%NUMJUSTPARAMS))
	}

	switch ind {
	case gr_slatAdvX:
		return int32(sl.advance.x)
	case gr_slatAdvY:
		return int32(sl.advance.y)
	case gr_slatAttTo:
		return boolToInt(sl.parent != nil)
	case gr_slatAttX:
		return int32(sl.attach.x)
	case gr_slatAttY:
		return int32(sl.attach.y)
	case gr_slatAttXOff, gr_slatAttYOff:
		return 0
	case gr_slatAttWithX:
		return int32(sl.with.x)
	case gr_slatAttWithY:
		return int32(sl.with.y)
	case gr_slatAttWithXOff, gr_slatAttWithYOff:
		return 0
	case gr_slatAttLevel:
		return int32(sl.attLevel)
	case gr_slatBreak:
		return int32(seg.getCharInfo(sl.original).breakWeight)
	case gr_slatCompRef:
		return 0
	case gr_slatDir:
		return int32(seg.dir & 1)
	case gr_slatInsert:
		return boolToInt(sl.isInsertBefore())
	case gr_slatPosX:
		return int32(sl.Position.x) // but need to calculate it
	case gr_slatPosY:
		return int32(sl.Position.y)
	case gr_slatShiftX:
		return int32(sl.shift.x)
	case gr_slatShiftY:
		return int32(sl.shift.y)
	case gr_slatMeasureSol:
		return -1 // err what's this?
	case gr_slatMeasureEol:
		return -1
	case gr_slatJWidth:
		return int32(sl.just)
	case gr_slatUserDefnV1:
		subindex = 0
		fallthrough
	case gr_slatUserDefn:
		if subindex < int(seg.silf.userAttibutes) {
			return int32(sl.userAttrs[subindex])
		}
	case gr_slatSegSplit:
		return int32(seg.getCharInfo(sl.original).flags & 3)
	case gr_slatBidiLevel:
		return int32(sl.bidiLevel)
	case gr_slatColFlags:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.flags)
		}
	case gr_slatColLimitblx:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.limit.bl.x)
		}
	case gr_slatColLimitbly:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.limit.bl.y)
		}
	case gr_slatColLimittrx:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.limit.tr.x)
		}
	case gr_slatColLimittry:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.limit.tr.y)
		}
	case gr_slatColShiftx:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.offset.x)
		}
	case gr_slatColShifty:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.offset.y)
		}
	case gr_slatColMargin:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.margin)
		}
	case gr_slatColMarginWt:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.marginWt)
		}
	case gr_slatColExclGlyph:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.exclGlyph)
		}
	case gr_slatColExclOffx:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.exclOffset.x)
		}
	case gr_slatColExclOffy:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.exclOffset.y)
		}
	case gr_slatSeqClass:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqClass)
		}
	case gr_slatSeqProxClass:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqProxClass)
		}
	case gr_slatSeqOrder:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqOrder)
		}
	case gr_slatSeqAboveXoff:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqAboveXoff)
		}
	case gr_slatSeqAboveWt:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqAboveWt)
		}
	case gr_slatSeqBelowXlim:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqBelowXlim)
		}
	case gr_slatSeqBelowWt:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqBelowWt)
		}
	case gr_slatSeqValignHt:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqValignHt)
		}
	case gr_slatSeqValignWt:
		if c := seg.getCollisionInfo(sl); c != nil {
			return int32(c.seqValignWt)
		}
	}
	return 0
}

func (sl *Slot) setJustify(seg *Segment, level uint8, subindex int, value int16) {
	if level != 0 && int(level) >= len(seg.silf.justificationLevels) {
		return
	}
	if sl.justs == nil {
		j := seg.newJustify()
		if j == nil {
			return
		}
		j.loadSlot(sl, seg)
		sl.justs = j
	}
	sl.justs.values[level][subindex] = value
}

func (sl *Slot) setAttr(map_ *slotMap, ind attrCode, subindex int, value int16) {
	seg := map_.segment
	if ind == gr_slatUserDefnV1 {
		ind = gr_slatUserDefn
		subindex = 0
		if seg.silf.userAttibutes == 0 {
			return
		}
	} else if ind >= gr_slatJStretch && ind < gr_slatJStretch+20 && ind != gr_slatJWidth {
		indx := int(ind - gr_slatJStretch)
		sl.setJustify(seg, uint8(indx/NUMJUSTPARAMS), indx%NUMJUSTPARAMS, value)
		return
	}

	switch ind {
	case gr_slatAdvX:
		sl.advance.x = float32(value)
	case gr_slatAdvY:
		sl.advance.y = float32(value)
	case gr_slatAttTo:
		idx := int(uint16(value))
		if idx < map_.size && map_.get(idx) != nil {
			other := map_.get(idx)
			if other == sl || other == sl.parent || other.isCopied() {
				break
			}
			if sl.parent != nil {
				sl.parent.removeChild(sl)
				sl.parent = nil
			}
			pOther := other
			count := 0
			foundOther := false
			for pOther != nil {
				count++
				if pOther == sl {
					foundOther = true
				}
				pOther = pOther.parent
			}
			for pOther = sl.child; pOther != nil; pOther = pOther.child {
				count++
			}
			for pOther = sl.sibling; pOther != nil; pOther = pOther.sibling {
				count++
			}
			if count < 100 && !foundOther && other.setChild(sl) {
				sl.parent = other
				if map_.dir != (idx > subindex) {
					sl.with = Position{sl.advance.x, 0}
				} else { // normal match to previous root
					sl.attach = Position{other.advance.x, 0}
				}
			}
		}
	case gr_slatAttX:
		sl.attach.x = float32(value)
	case gr_slatAttY:
		sl.attach.y = float32(value)
	case gr_slatAttXOff, gr_slatAttYOff:
	case gr_slatAttWithX:
		sl.with.x = float32(value)
	case gr_slatAttWithY:
		sl.with.y = float32(value)
	case gr_slatAttWithXOff, gr_slatAttWithYOff:
	case gr_slatAttLevel:
		sl.attLevel = byte(value)
	case gr_slatBreak:
		seg.getCharInfo(sl.original).breakWeight = value
	case gr_slatCompRef:
		// not sure what to do here
	case gr_slatDir:
	case gr_slatInsert:
		sl.markInsertBefore(value != 0)
	case gr_slatPosX:
		// can't set these here
	case gr_slatPosY:
	case gr_slatShiftX:
		sl.shift.x = float32(value)
	case gr_slatShiftY:
		sl.shift.y = float32(value)
	case gr_slatMeasureSol, gr_slatMeasureEol:
	case gr_slatJWidth:
		sl.just = float32(value)
	case gr_slatSegSplit:
		seg.getCharInfo(sl.original).addFlags(uint8(value & 3))
	case gr_slatUserDefn:
		sl.userAttrs[subindex] = value
	case gr_slatColFlags:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.flags = uint16(value)
		}
	case gr_slatColLimitblx:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			s := c.limit
			c.limit = rect{Position{float32(value), s.bl.y}, s.tr}
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColLimitbly:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			s := c.limit
			c.limit = rect{Position{s.bl.x, float32(value)}, s.tr}
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColLimittrx:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			s := c.limit
			c.limit = rect{s.bl, Position{float32(value), s.tr.y}}
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColLimittry:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			s := c.limit
			c.limit = rect{s.bl, Position{s.tr.x, float32(value)}}
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColMargin:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.margin = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColMarginWt:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.marginWt = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColExclGlyph:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.exclGlyph = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColExclOffx:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			s := c.exclOffset
			c.exclOffset = Position{float32(value), s.y}
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatColExclOffy:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			s := c.exclOffset
			c.exclOffset = Position{s.x, float32(value)}
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqClass:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqClass = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqProxClass:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqProxClass = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqOrder:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqOrder = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqAboveXoff:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqAboveXoff = value
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqAboveWt:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqAboveWt = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqBelowXlim:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqBelowXlim = value
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqBelowWt:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqBelowWt = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqValignHt:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqValignHt = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	case gr_slatSeqValignWt:
		c := seg.getCollisionInfo(sl)
		if c != nil {
			c.seqValignWt = uint16(value)
			c.flags = c.flags & ^COLL_KNOWN
		}
	}
}

func (sl *Slot) finalise(seg *Segment, font *graphiteFont, base Position, bbox *rect, attrLevel uint8, clusterMin *float32, rtl, isFinal bool, depth int) Position {
	if depth > 100 || (attrLevel != 0 && sl.attLevel > attrLevel) {
		return Position{}
	}
	var scale float32 = 1

	shift := Position{sl.shift.x*(float32(boolToInt(rtl)*-2+1)) + sl.just, sl.shift.y}
	tAdvance := sl.advance.x + sl.just
	if coll := seg.getCollisionInfo(sl); isFinal && coll != nil {
		collshift := coll.offset
		if coll.flags&COLL_KERN == 0 || rtl {
			shift = shift.add(collshift)
		}
	}
	glyphFace := seg.face.getGlyph(sl.glyphID)
	if font != nil {
		scale = font.scale
		shift = shift.scale(scale)
		if font.isHinted && glyphFace != nil {
			// tAdvance = (sl.advance.x-float32(glyphFace.advance)+sl.just)*scale + font.advance(sl.glyphID) // FIXME:
		} else {
			tAdvance *= scale
		}
	}
	var res Position

	sl.Position = base.add(shift)
	if sl.parent == nil {
		res = base.add(Position{tAdvance, sl.advance.y * scale})
		*clusterMin = sl.Position.x
	} else {
		sl.Position = sl.Position.add(sl.attach.sub(sl.with).scale(scale))
		var tAdv float32
		if sl.advance.x >= 0.5 {
			tAdv = sl.Position.x + tAdvance - shift.x
		}
		res = Position{tAdv, 0}
		if (sl.advance.x >= 0.5 || sl.Position.x < 0) && sl.Position.x < *clusterMin {
			*clusterMin = sl.Position.x
		}
	}

	if glyphFace != nil {
		ourBbox := glyphFace.bbox.scale(scale).addPosition(sl.Position)
		*bbox = bbox.widen(ourBbox)
	}

	if sl.child != nil && sl.child != sl && sl.child.parent == sl {
		tRes := sl.child.finalise(seg, font, sl.Position, bbox, attrLevel, clusterMin, rtl, isFinal, depth+1)
		if (sl.parent == nil || sl.advance.x >= 0.5) && tRes.x > res.x {
			res = tRes
		}
	}

	if sl.parent != nil && sl.sibling != nil && sl.sibling != sl && sl.sibling.parent == sl.parent {
		tRes := sl.sibling.finalise(seg, font, base, bbox, attrLevel, clusterMin, rtl, isFinal, depth+1)
		if tRes.x > res.x {
			res = tRes
		}
	}

	if sl.parent == nil && *clusterMin < base.x {
		adj := Position{sl.Position.x - *clusterMin, 0.}
		res = res.add(adj)
		sl.Position = sl.Position.add(adj)
		if sl.child != nil {
			sl.child.floodShift(adj, 0)
		}
	}
	return res
}

func (sl *Slot) floodShift(adj Position, depth int) {
	if depth > 100 {
		return
	}
	sl.Position = sl.Position.add(adj)
	if sl.child != nil {
		sl.child.floodShift(adj, depth+1)
	}
	if sl.sibling != nil {
		sl.sibling.floodShift(adj, depth+1)
	}
}

func (sl *Slot) clusterMetric(seg *Segment, metric, attrLevel uint8, rtl bool) int32 {
	if int(sl.glyphID) >= len(seg.face.glyphs) {
		return 0
	}
	bbox := seg.face.getGlyph(sl.glyphID).bbox
	var clusterMin float32

	res := sl.finalise(seg, nil, Position{}, &bbox, attrLevel, &clusterMin, rtl, false, 0)

	switch metric {
	case kgmetLsb:
		return int32(bbox.bl.x)
	case kgmetRsb:
		return int32(res.x - bbox.tr.x)
	case kgmetBbTop:
		return int32(bbox.tr.y)
	case kgmetBbBottom:
		return int32(bbox.bl.y)
	case kgmetBbLeft:
		return int32(bbox.bl.x)
	case kgmetBbRight:
		return int32(bbox.tr.x)
	case kgmetBbWidth:
		return int32(bbox.tr.x - bbox.bl.x)
	case kgmetBbHeight:
		return int32(bbox.tr.y - bbox.bl.y)
	case kgmetAdvWidth:
		return int32(res.x)
	case kgmetAdvHeight:
		return int32(res.y)
	default:
		return 0
	}
}

const NUMJUSTPARAMS = 5

type slotJustify struct {

	//     SlotJustify(const SlotJustify &);
	//     SlotJustify & operator = (const SlotJustify &);

	// public:
	//     static size_t size_of(size_t levels) { return sizeof(SlotJustify) + ((levels > 1 ? levels : 1)*NUMJUSTPARAMS - 1)*sizeof(int16); }

	//     void LoadSlot(const Slot *s, const Segment *seg);

	next   *slotJustify
	values [][NUMJUSTPARAMS]int16 // with length levels
}

func (sj *slotJustify) loadSlot(s *Slot, seg *Segment) {
	sj.values = make([][NUMJUSTPARAMS]int16, len(seg.silf.justificationLevels))
	for i, justs := range seg.silf.justificationLevels {
		v := &sj.values[i]
		v[0] = seg.face.getGlyphAttr(s.glyphID, uint16(justs.attrStretch))
		v[1] = seg.face.getGlyphAttr(s.glyphID, uint16(justs.attrShrink))
		v[2] = seg.face.getGlyphAttr(s.glyphID, uint16(justs.attrStep))
		v[3] = seg.face.getGlyphAttr(s.glyphID, uint16(justs.attrWeight))
	}
}