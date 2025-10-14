import json
import sys
import yaml

SWAGGER_UI_TEMPLATE = """
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>

<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
    <script>
        window.onload = () => {
            window.ui = SwaggerUIBundle({
                spec: {{ OPEN_API_SPEC }},
                dom_id: '#swagger-ui',
            });
        };
    </script>
</body>

</html>
"""


def generate_swagger_html(yaml_file, output_file):
    """Convert YAML OpenAPI spec to HTML with embedded Swagger UI"""
    with open(yaml_file, "r") as yaml_file_obj:
        openapi_spec = yaml.safe_load(yaml_file_obj)

    # Convert spec to JSON string for embedding
    spec_json = json.dumps(openapi_spec, indent=2)

    # Generate HTML with embedded spec
    html_content = SWAGGER_UI_TEMPLATE.replace("{{ OPEN_API_SPEC }}", spec_json)

    with open(output_file, "w") as file:
        file.write(html_content)

    print(f"Generated Swagger UI HTML at {output_file}")


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python genswagger-docs.py <openapi_yaml> <output_html>")
        sys.exit(1)

    openapi_yaml = sys.argv[1]
    output_html = sys.argv[2]

    generate_swagger_html(openapi_yaml, output_html)
