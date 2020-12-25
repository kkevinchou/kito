#version 330 core
out vec4 FragColor;

in vec3 Normal;  
in vec3 FragPos;
in vec2 TexCoord;
in vec4 Color;

// uniform vec3 lightPos; 
// uniform vec3 lightColor;
uniform vec3 viewPos;
uniform sampler2D ourTexture;

void main()
{
    vec3 lightPos = vec3(0, 20.0, 20.0);
    vec3 lightColor = vec3(1.0, 1.0, 1.0);

    // ambient
    float ambientStrength = 0.2;
    vec3 ambient = ambientStrength * lightColor;
        
    // diffuse 
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;

    FragColor = vec4((ambient + diffuse) * texture(ourTexture, TexCoord).xyz, 1.0);
    // FragColor = Color;
}