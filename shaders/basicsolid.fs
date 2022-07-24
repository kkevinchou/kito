#version 330 core
out vec4 FragColor;

in vec3 FragPos;  
in vec3 Normal;  
in float Alpha;

void main()
{
    vec3 objectColor = vec3(1.0, .0, .0);
    FragColor = vec4(objectColor, Alpha);
}
