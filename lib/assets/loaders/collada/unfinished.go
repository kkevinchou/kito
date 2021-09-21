package collada

//AnimationClip defines a section of the animation curves to be used together as an animation clip.
type AnimationClip struct {
	//TODO
}

//InstanceAnimation instantiates a COLLADA animation resource.
type InstanceAnimation struct {
	//TODO
}

//LibraryAnimationClips provides a library in which to place <animation_clip> elements.
type LibraryAnimationClips struct {
	//TODO
}

//Orthographic describes the field of view of an orthographic camera.
type Orthographic struct {
	//TODO
}

//Perspective describes the field of view of a perspective camera.
type Perspective struct {
	//TODO
}

//Morph describes the data required to blend between sets of static meshes.
type Morph struct {
	//TODO
}

//Skeleton indicates where a skin controller is to start searching for the joint nodes that it needs.
type Skeleton struct {
	//TODO
}

// Targets teclares morph targets, their weights, and any user-defined attributes associated with them.
type Targets struct {
	//TODO
}

//AmbientCore (core) Describes an ambient light source.
type AmbientCore struct {
	//TODO
}

//Directional describes a directional light source.
type Directional struct {
	//TODO
}

//Point describes a point light source.
type Point struct {
	//TODO
}

//Spot describes a spot light source.
type Spot struct {
	//TODO
}

//Formula defines a formula.
type Formula struct {
	//TODO
}

//InstanceFormula instantiates a COLLADA formula resource.
type InstanceFormula struct {
	//TODO
}

//LibraryFormulas provides a library in which to place <formula> elements.
type LibraryFormulas struct {
	//TODO
}

//GeographicLocation defines an asset’s location for asset management.
type GeographicLocation struct {
	//TODO
}

//EvaluateScene declares information specifying how to evaluate a <visual_scene>.
type EvaluateScene struct {
	//TODO
}

//LibraryNodes provides a library in which to place <node> elements.
type LibraryNodes struct {
	//TODO
}

type InstancePhysicsScene struct {
	//TODO
}

type InstanceKinematicsScene struct {
	//TODO
}

//Array Creates a parameter of a one-dimensional array type.
type Array struct {
	//TODO
}

//Modifier Provides additional information about the volatility or linkage of a <newparam>declaration.
type Modifier struct {
	//TODO
}

//Newparam Creates a new, named parameter object and assigns it a type and an initial value. See Chapter 5: Core Elements Reference.
type Newparam struct {
	//TODO
}

//ParamReference (reference) References a predefined parameter. See Chapter 5: Core Elements Reference.
type ParamReference struct {
	//TODO
}

//SamplerStates Allows users to modify an effect’s sampler state from a material.
type SamplerStates struct {
	//TODO
}

//Semantic Provides metadata that describes the purpose of a parameter declaration.
type Semantic struct {
	//TODO
}

//Usertype Creates an instance of a structured class for a parameter.
type Usertype struct {
	//TODO
}

//ProfileCg Declares a platform-specific representation of an effect written in the NVIDIA®Cg language.
type ProfileCg struct {
	//TODO
}
type ProfileGles struct {
	//TODO
}
type ProfileGles2 struct {
	//TODO
}
type ProfileGlsl struct {
	//TODO
}
type Blinn struct {
	//TODO
}
type ColorClear struct {
	//TODO
}
type ColorTarget struct {
	//TODO
}

//Constant Produces a constantly shaded surface that is independent of lighting.
type ConstantFx struct {
	//TODO
}

//DepthClear Specifies whether a render target image is to be cleared, and which value to use.
type DepthClear struct {
	//TODO
}

//DepthTarget Specifies which <image> will receive the depth information from the output of this pass.
type DepthTarget struct {
	//TODO
}

//Draw Instructs the FX Runtime what kind of geometry to submit.
type Draw struct {
	//TODO
}

//Evaluate Contains evaluation elements for a rendering pass.
type Evaluate struct {
	//TODO
}

//InstanceMaterialRendering Instantiates a COLLADA material resource for a screen effect.
type InstanceMaterialRendering struct {
	//TODO
}

//Lambert Produces a diffuse shaded surface that is independent of lighting.
type Lambert struct {
	//TODO
}

//Pass Provides a static declaration of all the render states, shaders, and settings for one rendering pipeline.
type Pass struct {
	//TODO
}

//According the Phong BRDF approximation.
type According struct {
	//TODO
}

//Render Describes one effect pass to evaluate a scene.
type Render struct {
	//TODO
}

//States Contains all rendering states to set up for the parent pass.
type States struct {
	//TODO
}

//StencilClear Specifies whether a render target image is to be cleared, and which value to use.
type StencilClear struct {
	//TODO
}

//StencilTarget Specifies which <image> will receive the stencil information from the output of this pass
type StencilTarget struct {
	//TODO
}

//Binary Identifies or provides a shader in binary form.
type Binary struct {
	//TODO
}

//BindAttribute Binds semantics to vertex attribute inputs of a shader.
type BindAttribute struct {
	//TODO
}

//BindUniform Binds values to uniform inputs of a shader or binds values to effect parameters upon instantiation.
type BindUniform struct {
	//TODO
}

//Code Provides an inline block of source code.
type Code struct {
	//TODO
}

//Compiler Contains command-line or runtime-invocation options for a shader compiler.
type Compiler struct {
	//TODO
}

//Include Imports source code or precompiled binary shaders into the FX Runtime by referencing an external resource.
type Include struct {
	//TODO
}

//Linker Contains command-line or runtime-invocation options for shader linkers to combine shaders into programs.
type Linker struct {
	//TODO
}

//Program Links multiple shaders together to produce a pipeline for geometry processing.
type Program struct {
	//TODO
}
type Shader struct {
	//TODO
}
type Sources struct {
	//TODO
}
type Alpha struct {
	//TODO
}
type Argument struct {
	//TODO
}
type Create2d struct {
	//TODO
}
type Create3d struct {
	//TODO
}
type CreateCube struct {
	//TODO
}
type Format struct {
	//TODO
}
type Image struct {
	//TODO
}
type InitFrom struct {
	//TODO
}
type InstanceImage struct {
	//TODO
}
type LibraryImages struct {
	//TODO
}
type Rgb struct {
	//TODO
}
type FxSamplerCommon struct {
	//TODO
}
type Sampler1D struct {
	//TODO
}
type Sampler2D struct {
	//TODO
}
type Sampler3D struct {
	//TODO
}
type SamplerCube struct {
	//TODO
}
type SamplerDepth struct {
	//TODO
}
type SamplerRect struct {
	//TODO
}
type Texcombiner struct {
	//TODO
}
type Texenv struct {
	//TODO
}

//TexturePipeline Defines a set of texturing commands that will be converted into multitexturing operations using glTexEnv in regular and combiner mode.
type TexturePipeline struct {
	//TODO
}
