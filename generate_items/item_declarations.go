package main

import (
  "github.com/wowsims/tbc/generate_items/api"
)

type ItemDeclaration struct {
  ID int
  Specs []api.Spec // Which specs use this item
}

// Keep these sorted by ID.
var ItemDeclarations = []ItemDeclaration{
	{ ID: 19344, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 19379, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 21608, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 21709, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 22730, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23025, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23031, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23046, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23049, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23050, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23057, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23070, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23199, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23207, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23554, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23664, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 23665, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24116, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24121, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24126, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24250, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24252, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24256, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24262, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24266, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24452, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 24557, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 25777, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 25778, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27462, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27464, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27465, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27469, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27470, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27471, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27472, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27473, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27488, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27492, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27493, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27508, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27510, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27522, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27534, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27537, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27543, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27683, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27741, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27743, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27746, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27758, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27778, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27783, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27784, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27793, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27795, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27796, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27802, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27821, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27824, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27838, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27842, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27845, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27868, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27907, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27909, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27910, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27914, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27937, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27948, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27981, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27993, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 27994, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28134, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28169, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28174, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28179, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28185, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28187, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28188, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28191, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28193, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28227, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28229, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28231, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28232, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28245, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28248, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28254, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28260, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28266, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28269, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28278, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28297, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28341, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28342, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28346, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28349, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28379, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28391, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28394, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28406, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28412, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28415, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28418, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28507, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28510, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28515, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28517, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28530, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28555, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28565, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28570, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28583, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28585, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28586, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28594, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28602, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28603, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28611, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28633, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28638, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28639, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28640, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28654, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28670, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28726, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28734, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28744, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28753, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28758, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28762, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28766, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28770, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28780, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28781, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28785, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28789, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28793, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28797, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28799, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 28810, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29033, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29034, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29035, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29036, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29037, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29126, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29129, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29130, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29132, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29141, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29142, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29172, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29179, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29240, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29241, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29242, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29243, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29244, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29245, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29257, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29258, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29268, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29273, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29285, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29286, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29287, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29302, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29305, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29313, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29314, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29317, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29320, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29330, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29333, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29341, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29343, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29352, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29355, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29367, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29368, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29369, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29370, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29376, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29389, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29504, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29519, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29520, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29521, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29522, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29523, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29524, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29784, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29808, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29813, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29918, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29922, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29955, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29972, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29986, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29987, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29988, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 29992, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30004, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30008, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30011, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30015, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30024, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30037, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30038, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30043, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30044, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30049, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30056, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30064, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30067, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30079, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30107, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30109, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30169, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30170, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30171, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30172, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30173, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30297, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30366, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30519, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30531, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30532, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30541, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30626, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30663, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30667, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30677, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30682, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30686, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30709, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30723, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30725, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30734, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30735, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30832, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30870, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30872, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30884, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30888, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30894, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30909, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30913, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30914, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30916, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30924, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30925, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30946, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 30984, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31008, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31014, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31017, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31020, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31023, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31075, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31107, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31140, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31149, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31280, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31283, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31287, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31290, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31297, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31308, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31330, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31338, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31339, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31340, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31461, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31513, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31692, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31693, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31797, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31856, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31921, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 31922, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32078, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32086, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32237, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32239, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32242, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32247, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32256, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32259, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32270, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32276, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32327, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32328, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32330, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32331, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32338, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32349, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32351, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32352, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32361, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32367, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32374, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32476, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32483, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32524, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32525, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32527, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32541, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32586, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32587, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32592, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32664, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32779, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32792, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32817, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 32963, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33281, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33283, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33334, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33354, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33357, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33466, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33506, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33533, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33534, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33537, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33588, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33591, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33829, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33965, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 33970, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34009, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34011, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34179, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34186, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34204, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34230, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34242, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34332, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34336, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34344, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34350, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34359, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34362, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34390, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34396, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34429, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34437, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34542, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 34566, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 35749, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
	{ ID: 38290, Specs: []api.Spec{api.Spec_SpecElementalShaman}, },
}
