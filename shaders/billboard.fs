#version 330 core
out vec4 FragColor;

in vec3 FragPos;  

void main()
{
    vec3 objectColor = vec3(1.0, .0, .0);
    FragColor = vec4(objectColor, 1.0);
}
