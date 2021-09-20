#version 330 core
const int MAX_JOINTS = 50;
const int MAX_WEIGHTS = 3;

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;
layout (location = 3) in vec3 aColor;
layout (location = 4) in ivec3 jointIndices;
layout (location = 5) in vec3 jointWeights;

out VS_OUT {
    vec3 FragPos;
    vec3 Normal;
    vec4 FragPosLightSpace;
    mat4 View;
    vec2 TexCoord;
    flat int SpecialVert;
    flat int SpecialVert2;
    flat int SpecialVert3;
} vs_out;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat4 jointTransforms[MAX_JOINTS];
uniform mat4 lightSpaceMatrix;

void main() {
    vec4 totalLocalPos = vec4(0.0);
	vec4 totalNormal = vec4(0.0);

    int specialVert = 0;
    int specialVert2 = 0;

    float sumWeights = 0;
	for(int i = 0; i < MAX_WEIGHTS; i++){
		int jointIndex = jointIndices[i];
        // if (jointIndex == 50 || jointIndex == 51) {
        //     specialVert = 1;
        // }

        // if (jointIndex == 38 && jointWeights[i] >= 0.7) {
        //     specialVert2 = 1;
        // }

        if (jointIndex == 50) {
            specialVert2 = 1;
        }

		mat4 jointTransform = jointTransforms[jointIndex];
		vec4 posePosition = jointTransform * vec4(aPos, 1.0);
		totalLocalPos += posePosition * jointWeights[i];

		vec4 worldNormal = jointTransform * vec4(aNormal, 0.0);
		totalNormal += worldNormal * jointWeights[i];
        sumWeights += jointWeights[i];
	}

    vs_out.SpecialVert = specialVert;
    vs_out.SpecialVert2 = specialVert2;
    if (sumWeights <= 0.99) {
        vs_out.SpecialVert3 = 1;
    }


    vs_out.FragPos = vec3(model * totalLocalPos);
    // TODO: the normal matrix is expensive to calculate and should be passed in as a uniform
    vs_out.Normal = transpose(inverse(mat3(model))) * vec3(totalNormal);
    vs_out.FragPosLightSpace = lightSpaceMatrix * vec4(vs_out.FragPos, 1.0);
    vs_out.View = view;
	vs_out.TexCoord = aTexCoord;

    gl_Position = (projection * (view * (model * totalLocalPos)));
}