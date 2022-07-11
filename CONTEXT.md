# Context

## GLTF Animation exporting
When exporting to GLTF format, the blender exporter has a default on setting called "Optimize Animation Size". Which has a description of:
"Reduces exported filesize by removing duplicate keyframes. Can cause problems with stepped animation". Our current animation player can
only handle animated models that have keyframes for all joints (and probably for all frames as well). So, to fix my broken animations i
had to uncheck that setting. With it checked, the exported animations would have the first keyframe having transforms for all joints, while
subsequent keyframes would only have transforms for a portion of the joints. Presumably because it dropped duplicate keyframes. The GLTF spec
probably describes how to properly handle these scenarios (the vscode model renderer handles this properly). So this will be something for me
to implement in the GLTF loader or animation player. My exporting steps are to select nothing and export with "Optimize Animation Size" disabled.
It looks like this feature was added in for newer versions of blender, but when I was doing development this option was not available. Each
animation should be stashed in the same NLA track (i think)