package structs

const XahttpOpt_T = 1
const XahttpOpt_F = 0

var XahttpOpt = func(i, u, r, f, s int) byte {
	return byte(i&1<<3 + u&1<<2 + r&1<<1 + f&1 + s&1<<4)
}

var XahttpOptNoS = func() byte {
	return XahttpOpt(XahttpOpt_T, XahttpOpt_T, XahttpOpt_T, XahttpOpt_T, XahttpOpt_F)
}

var XahttpOptNoU = func() byte {
	return XahttpOpt(XahttpOpt_T, XahttpOpt_F, XahttpOpt_T, XahttpOpt_T, XahttpOpt_F)
}

var XahttpOptNoR = func() byte {
	return XahttpOpt(XahttpOpt_T, XahttpOpt_T, XahttpOpt_F, XahttpOpt_T, XahttpOpt_F)
}

var XahttpOptNoUR = func() byte {
	return XahttpOpt(XahttpOpt_T, XahttpOpt_F, XahttpOpt_F, XahttpOpt_T, XahttpOpt_F)
}

var XahttpOptIO = func() byte {
	return XahttpOpt(XahttpOpt_T, XahttpOpt_F, XahttpOpt_F, XahttpOpt_F, XahttpOpt_F)
}

var XahttpOptFO = func() byte {
	return XahttpOpt(XahttpOpt_F, XahttpOpt_F, XahttpOpt_F, XahttpOpt_T, XahttpOpt_F)
}

var XahttpOptSO = func() byte {
	return XahttpOpt(XahttpOpt_F, XahttpOpt_F, XahttpOpt_F, XahttpOpt_F, XahttpOpt_T)
}

var XahttpOptAll = func() byte {
	return XahttpOpt(XahttpOpt_T, XahttpOpt_T, XahttpOpt_T, XahttpOpt_T, XahttpOpt_T)
}

var XsacPhysically = func() byte {
	return 0b10000000
}

var XsacTombstone = func() byte {
	return 0
}

var XsacTombstoneAndHistory = func() byte {
	return 0b00000001
}

var XsacTombstoneAndForce = func() byte {
	return 0b00000011
}

var XsacTombstoneAndRestore = func() byte {
	return 0b00000101
}

var XsacTombstoneWhole = func() byte {
	return 0b00000111
}
