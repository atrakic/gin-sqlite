import yaml, sys
import json

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
                url: '{{ OPEN_API_URL }}',
                dom_id: '#swagger-ui',
            });
        };
    </script>
</body>

</html>
"""


def generate_swagger_ui(url, output_file):
    with open(output_file, "w") as file:
        file.write(SWAGGER_UI_TEMPLATE.replace("{{ OPEN_API_URL }}", url))
    print(f"Generated Swagger UI at {output_file}")


def convert_yaml_to_json(yaml_file, output_file):
    with open(yaml_file, "r") as yaml_file_obj:
        with open(output_file, "w") as file:
            json.dump(yaml.safe_load(yaml_file_obj), file, indent=2)
    print(f"Converted {yaml_file} to {output_file}")


if __name__ == "__main__":
    if len(sys.argv) != 4:
        print(
            "Usage: python genswagger-docs.py <openapi_yaml> <output_json> <output_html>"
        )
        sys.exit(1)

    openapi_yaml = sys.argv[1]
    output_json = sys.argv[2]
    output_html = sys.argv[3]

    convert_yaml_to_json(openapi_yaml, output_json)
    generate_swagger_ui(output_json, output_html)
