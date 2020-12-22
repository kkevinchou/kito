package collada

//Skin contains vertex and primitive information sufficient to describe blend-weight skinning.
type Skin struct {
	Source        []*Source     `xml:"source"`
	VertexWeights VertexWeights `xml:"vertex_weights"`
}

// VertexWeights describes the combination of joints and weights used by a skin.
type VertexWeights struct {
	VCount string `xml:"vcount"`
	V      string `xml:"v"`
}

//Controller categorizes the declaration of generic control information.
type Controller struct {
	Skin Skin `xml:"skin"`
}

//LibraryControllers provides a library in which to place <controller> elements.
type LibraryControllers struct {
	Controller Controller `xml:"controller"`
}

type HasTechniqueCommon struct {
	TechniqueCommon TechniqueCommon `xml:"technique_common"`
}

//TechniqueCommon specifies the information for a specific element for the common profile that all COLLADA implementations must support.
type TechniqueCommon struct {
	// XML string `xml:",innerxml"`
	Accessor *Accessor `xml:"accessor"`
}

// Accessor declares an access pattern to one of the array elements <float_array>, <int_array>, <Name_array>, <bool_array>, and <IDREF_array>.
type Param struct {
	Name string `xml:"name,attr,omitempty"`
}

type Accessor struct {
	Param *Param `xml:"param"`
}
