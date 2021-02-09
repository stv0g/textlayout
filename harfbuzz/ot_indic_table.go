package harfbuzz

// Code generated by unicodedata/generate/main.go DO NOT EDIT.

const (
	OT_X            = 0
	OT_C            = 1
	OT_V            = 2
	OT_N            = 3
	OT_H            = 4
	OT_ZWNJ         = 5
	OT_ZWJ          = 6
	OT_M            = 7
	OT_SM           = 8
	OT_A            = 10
	OT_PLACEHOLDER  = 11
	OT_DOTTEDCIRCLE = 12
	OT_RS           = 13
	OT_Coeng        = 14
	OT_Repha        = 15
	OT_Ra           = 16
	OT_CM           = 17
	OT_Symbol       = 18
	OT_CS           = 19
	OT_VAbv         = 26
	OT_VBlw         = 27
	OT_VPre         = 28
	OT_VPst         = 29
)
const (
	POS_START             = 0
	POS_RA_TO_BECOME_REPH = 1
	POS_PRE_M             = 2
	POS_PRE_C             = 3
	POS_BASE_C            = 4
	POS_AFTER_MAIN        = 5
	POS_ABOVE_C           = 6
	POS_BEFORE_SUB        = 7
	POS_BELOW_C           = 8
	POS_AFTER_SUB         = 9
	POS_BEFORE_POST       = 10
	POS_POST_C            = 11
	POS_AFTER_POST        = 12
	POS_FINAL_C           = 13
	POS_SMVD              = 14
	POS_END               = 15
)

var indicTable = [...]uint16{

	/* Basic Latin */
	0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf0b, 0xf00, 0xf00,
	/* 0030 */ 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,

	/* Latin-1 Supplement */

	/* 00B0 */ 0xf00, 0xf00, 0xf08, 0xf08, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,
	/* 00C0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,
	/* 00D0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf0b,

	/* Devanagari */

	/* 0900 */ 0xf08, 0xf08, 0xf08, 0xf08, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02,
	/* 0910 */ 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0920 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0930 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0x607, 0xb07, 0xf03, 0xf12, 0xb07, 0x307,
	/* 0940 */ 0xb07, 0x807, 0x807, 0x807, 0x807, 0x607, 0x607, 0x607, 0x607, 0xb07, 0xb07, 0xb07, 0xb07, 0x804, 0x307, 0xb07,
	/* 0950 */ 0xf00, 0xf0a, 0xf0a, 0xf00, 0xf00, 0x607, 0x807, 0x807, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0960 */ 0xf02, 0xf02, 0x807, 0x807, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0970 */ 0xf00, 0xf00, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,

	/* Bengali */

	/* 0980 */ 0xf0b, 0xf08, 0xf08, 0xf08, 0xf00, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf00, 0xf02,
	/* 0990 */ 0xf02, 0xf00, 0xf00, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 09A0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 09B0 */ 0xf01, 0xf00, 0xf01, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf03, 0xf12, 0xb07, 0x307,
	/* 09C0 */ 0xb07, 0x807, 0x807, 0x807, 0x807, 0xf00, 0xf00, 0x307, 0x307, 0xf00, 0xf00, 0xb07, 0xb07, 0x804, 0xf01, 0xf00,
	/* 09D0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xb07, 0xf00, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf00, 0xf01,
	/* 09E0 */ 0xf02, 0xf02, 0x807, 0x807, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 09F0 */ 0xf01, 0xf01, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf08, 0xf00, 0xf08, 0xf00,

	/* Gurmukhi */

	/* 0A00 */ 0xf00, 0xf08, 0xf08, 0xf08, 0xf00, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf00, 0xf00, 0xf00, 0xf02,
	/* 0A10 */ 0xf02, 0xf00, 0xf00, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0A20 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0A30 */ 0xf01, 0xf00, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf00, 0xf00, 0xf03, 0xf00, 0xb07, 0x307,
	/* 0A40 */ 0xb07, 0x807, 0x807, 0xf00, 0xf00, 0xf00, 0xf00, 0x607, 0x607, 0xf00, 0xf00, 0x607, 0x607, 0x804, 0xf00, 0xf00,
	/* 0A50 */ 0xf00, 0xf0a, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf00,
	/* 0A60 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0A70 */ 0xf08, 0x608, 0xf0b, 0xf0b, 0xf00, 0x811, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,

	/* Gujarati */

	/* 0A80 */ 0xf00, 0xf08, 0xf08, 0xf08, 0xf00, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf02,
	/* 0A90 */ 0xf02, 0xf02, 0xf00, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0AA0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0AB0 */ 0xf01, 0xf00, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf03, 0xf12, 0xb07, 0x307,
	/* 0AC0 */ 0xb07, 0x807, 0x807, 0x807, 0x807, 0x607, 0xf00, 0x607, 0x607, 0xb07, 0xf00, 0xb07, 0xb07, 0x804, 0xf00, 0xf00,
	/* 0AD0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,
	/* 0AE0 */ 0xf02, 0xf02, 0x807, 0x807, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0AF0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf01, 0xf0a, 0xf0a, 0xf0a, 0xf03, 0xf03, 0xf03,
	/* 0B00 */ 0xf00, 0xf08, 0xf08, 0xf08, 0xf00, 0xf02, 0xf02, 0xf02,

	/* Oriya */
	0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf00, 0xf02,
	/* 0B10 */ 0xf02, 0xf00, 0xf00, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0B20 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0B30 */ 0xf01, 0xf00, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf03, 0xf12, 0xb07, 0x607,
	/* 0B40 */ 0xb07, 0x807, 0x807, 0x807, 0x807, 0xf00, 0xf00, 0x307, 0x607, 0xf00, 0xf00, 0xb07, 0xb07, 0x804, 0xf00, 0xf00,
	/* 0B50 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0x607, 0x607, 0xb07, 0xf00, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf00, 0xf01,
	/* 0B60 */ 0xf02, 0xf02, 0x807, 0x807, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0B70 */ 0xf00, 0xf01, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,

	/* Tamil */

	/* 0B80 */ 0xf00, 0xf00, 0xf08, 0xf00, 0xf00, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf00, 0xf00, 0xf02, 0xf02,
	/* 0B90 */ 0xf02, 0xf00, 0xf02, 0xf02, 0xf02, 0xf01, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf00, 0xf01, 0xf00, 0xf01, 0xf01,
	/* 0BA0 */ 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01,
	/* 0BB0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf00, 0xf00, 0xb07, 0xb07,
	/* 0BC0 */ 0x607, 0xb07, 0xb07, 0xf00, 0xf00, 0xf00, 0x307, 0x307, 0x307, 0xf00, 0xb07, 0xb07, 0xb07, 0x604, 0xf00, 0xf00,
	/* 0BD0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xb07, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,
	/* 0BE0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0BF0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,

	/* Telugu */

	/* 0C00 */ 0xf08, 0xf08, 0xf08, 0xf08, 0xf08, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf02, 0xf02,
	/* 0C10 */ 0xf02, 0xf00, 0xf02, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0C20 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0C30 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf00, 0xf12, 0x607, 0x607,
	/* 0C40 */ 0x607, 0xb07, 0xb07, 0xb07, 0xb07, 0xf00, 0x607, 0x607, 0x807, 0xf00, 0x607, 0x607, 0x607, 0x604, 0xf00, 0xf00,
	/* 0C50 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0x607, 0x807, 0xf00, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,
	/* 0C60 */ 0xf02, 0xf02, 0x807, 0x807, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0C70 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,

	/* Kannada */

	/* 0C80 */ 0xf08, 0xf08, 0xf08, 0xf08, 0xf00, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf02, 0xf02,
	/* 0C90 */ 0xf02, 0xf00, 0xf02, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0CA0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0CB0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf03, 0xf12, 0xb07, 0x607,
	/* 0CC0 */ 0xb07, 0xb07, 0xb07, 0xb07, 0xb07, 0xf00, 0x607, 0xb07, 0xb07, 0xf00, 0xb07, 0xb07, 0x607, 0x604, 0xf00, 0xf00,
	/* 0CD0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xb07, 0xb07, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf01, 0xf00,
	/* 0CE0 */ 0xf02, 0xf02, 0x807, 0x807, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0CF0 */ 0xf00, 0xf13, 0xf13, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,

	/* Malayalam */

	/* 0D00 */ 0xf08, 0xf08, 0xf08, 0xf08, 0xf08, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf02, 0xf02,
	/* 0D10 */ 0xf02, 0xf00, 0xf02, 0xf02, 0xf02, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0D20 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0D30 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf07, 0xf07, 0xf12, 0xb07, 0xb07,
	/* 0D40 */ 0xb07, 0xb07, 0xb07, 0x807, 0x807, 0xf00, 0x307, 0x307, 0x307, 0xf00, 0xb07, 0xb07, 0xb07, 0x604, 0xf0f, 0xf00,
	/* 0D50 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf01, 0xb07, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf02,
	/* 0D60 */ 0xf02, 0xf02, 0x807, 0x807, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,
	/* 0D70 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0D80 */ 0xf00, 0xf08, 0xf08, 0xf08, 0xf00, 0xf02, 0xf02, 0xf02,

	/* Sinhala */
	0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02,
	/* 0D90 */ 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf00, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0DA0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 0DB0 */ 0xf01, 0xf01, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf01, 0xf00, 0xf00,
	/* 0DC0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf00, 0xf00, 0xf00, 0x604, 0xf00, 0xf00, 0xf00, 0xf00, 0xb07,
	/* 0DD0 */ 0xb07, 0xb07, 0x607, 0x607, 0x807, 0xf00, 0x807, 0xf00, 0xb07, 0x307, 0x607, 0x307, 0xb07, 0xb07, 0xb07, 0xb07,
	/* 0DE0 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,

	/* Myanmar */

	/* 1000 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 1010 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 1020 */ 0xf01, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xb07, 0xb07, 0x607, 0x607, 0x807,
	/* 1030 */ 0x807, 0x307, 0x607, 0x607, 0x607, 0x607, 0xf08, 0xf03, 0xf08, 0xf0e, 0xf07, 0xb11, 0x811, 0x811, 0x811, 0xf01,
	/* 1040 */ 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf00, 0xf0b, 0xf00, 0xf00, 0xf0b, 0xf00,
	/* 1050 */ 0xf01, 0xf01, 0xf02, 0xf02, 0xf02, 0xf02, 0xb07, 0xb07, 0x807, 0x807, 0xf01, 0xf01, 0xf01, 0xf01, 0x811, 0x811,
	/* 1060 */ 0x811, 0xf01, 0xb07, 0xf03, 0xf03, 0xf01, 0xf01, 0xb07, 0xb07, 0xf03, 0xf03, 0xf03, 0xf03, 0xf03, 0xf01, 0xf01,
	/* 1070 */ 0xf01, 0x607, 0x607, 0x607, 0x607, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 1080 */ 0xf01, 0xf01, 0x811, 0xb07, 0x307, 0x607, 0x607, 0xf03, 0xf03, 0xf03, 0xf03, 0xf03, 0xf03, 0xf03, 0xf01, 0xf03,
	/* 1090 */ 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,

	/* Khmer */

	/* 1780 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 1790 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* 17A0 */ 0xf01, 0xf01, 0xf01, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02, 0xf02,
	/* 17B0 */ 0xf02, 0xf02, 0xf02, 0xf02, 0xf00, 0xf00, 0xb07, 0x607, 0x607, 0x607, 0x607, 0x807, 0x807, 0x807, 0x607, 0xb07,
	/* 17C0 */ 0xb07, 0x307, 0x307, 0x307, 0xb07, 0xb07, 0xf08, 0xf08, 0xb07, 0x60d, 0x60d, 0xf08, 0x611, 0xf07, 0xf08, 0xf08,
	/* 17D0 */ 0xf08, 0xf07, 0xf0e, 0xf08, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf12, 0xf08, 0xf00, 0xf00,
	/* 17E0 */ 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,

	/* Vedic Extensions */

	/* 1CD0 */ 0xf0a, 0xf0a, 0xf0a, 0xf00, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a,
	/* 1CE0 */ 0xf0a, 0xf0a, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,
	/* 1CF0 */ 0xf00, 0xf00, 0xf01, 0xf01, 0xf0a, 0xf13, 0xf13, 0xf0a,

	/* General Punctuation */
	0xf00, 0xf00, 0xf00, 0xf00, 0xf05, 0xf06, 0xf00, 0xf00,

	/* Superscripts and Subscripts */

	/* 2070 */ 0xf00, 0xf00, 0xf00, 0xf00, 0xf08, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00, 0xf00,
	/* 2080 */ 0xf00, 0xf00, 0xf08, 0xf08, 0xf08, 0xf00, 0xf00, 0xf00,

	/* Devanagari Extended */

	/* A8E0 */ 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a, 0xf0a,
	/* A8F0 */ 0xf0a, 0xf0a, 0xf08, 0xf08, 0xf00, 0xf00, 0xf00, 0xf00,

	/* Myanmar Extended-B */

	/* A9E0 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0x607, 0xf00, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* A9F0 */ 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b, 0xf0b,

	/* Myanmar Extended-A */

	/* AA60 */ 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01, 0xf01,
	/* AA70 */ 0xf00, 0xf01, 0xf01, 0xf01, 0xf0b, 0xf0b, 0xf0b, 0xf00,
} /* Table items: 1728; occupancy: 71% */

const (
	offsetIndic0x0028u = 0
	offsetIndic0x00b0u = 24
	offsetIndic0x0900u = 64
	offsetIndic0x1000u = 1328
	offsetIndic0x1780u = 1480
	offsetIndic0x1cd0u = 1584
	offsetIndic0x2008u = 1624
	offsetIndic0x2070u = 1632
	offsetIndic0xa8e0u = 1656
	offsetIndic0xa9e0u = 1680
	offsetIndic0xaa60u = 1704
)

func indicGetCategories(u rune) uint16 {
	switch u >> 12 {
	case 0x0:
		if u == 0x00A0 {
			return 0xf0b
		}
		if 0x0028 <= u && u <= 0x003F {
			return indicTable[u-0x0028+offsetIndic0x0028u]
		}
		if 0x00B0 <= u && u <= 0x00D7 {
			return indicTable[u-0x00B0+offsetIndic0x00b0u]
		}
		if 0x0900 <= u && u <= 0x0DEF {
			return indicTable[u-0x0900+offsetIndic0x0900u]
		}

	case 0x1:
		if 0x1000 <= u && u <= 0x1097 {
			return indicTable[u-0x1000+offsetIndic0x1000u]
		}
		if 0x1780 <= u && u <= 0x17E7 {
			return indicTable[u-0x1780+offsetIndic0x1780u]
		}
		if 0x1CD0 <= u && u <= 0x1CF7 {
			return indicTable[u-0x1CD0+offsetIndic0x1cd0u]
		}

	case 0x2:
		if u == 0x25CC {
			return 0xf0b
		}
		if 0x2008 <= u && u <= 0x200F {
			return indicTable[u-0x2008+offsetIndic0x2008u]
		}
		if 0x2070 <= u && u <= 0x2087 {
			return indicTable[u-0x2070+offsetIndic0x2070u]
		}

	case 0xA:
		if 0xA8E0 <= u && u <= 0xA8F7 {
			return indicTable[u-0xA8E0+offsetIndic0xa8e0u]
		}
		if 0xA9E0 <= u && u <= 0xA9F7 {
			return indicTable[u-0xA9E0+offsetIndic0xa9e0u]
		}
		if 0xAA60 <= u && u <= 0xAA77 {
			return indicTable[u-0xAA60+offsetIndic0xaa60u]
		}

	}
	return 0xf00
}
