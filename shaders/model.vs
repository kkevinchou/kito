#version 330 core
const int MAX_JOINTS = 50;
const int MAX_WEIGHTS = 3;

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;
layout (location = 3) in vec3 aColor;
layout (location = 4) in ivec3 jointIndices;
layout (location = 5) in vec3 jointWeights;

out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat4 jointTransforms[MAX_JOINTS];

void main() {
    vec4 totalLocalPos = vec4(0.0);
	vec4 totalNormal = vec4(0.0);

	for(int i=0;i<MAX_WEIGHTS;i++){
		mat4 jointTransform = jointTransforms[jointIndices[i]];
		vec4 posePosition = jointTransform * vec4(aPos, 1.0);
		totalLocalPos += posePosition * jointWeights[i];

		vec4 worldNormal = jointTransform * vec4(aNormal, 0.0);
		totalNormal += worldNormal * jointWeights[i];
	}

	// totalLocalPos = vec4(aPos, 1);
	// totalNormal = vec4(aNormal, 1);

    FragPos = vec3(model * totalLocalPos);
    // TODO: the normal matrix is expensive to calculate and should be passed in as a uniform
    Normal = mat3(transpose(inverse(model))) * vec3(totalNormal);
    TexCoord = aTexCoord;

    gl_Position = (projection * (view * (model * totalLocalPos)));
    // gl_Position = projection * view * vec4(FragPos, 1.0);
}