import requests
import bs4
from golang_code_generator.golang_code_gen import (
    fields_to_struct,
    to_field,
    snake_case_to_pascal,
    Field,
    GoFile
)


def create_wprdc_struct():
    url = """https://data.wprdc.org/api/3/action/datastore_search_sql?sql=SELECT * FROM "518b583f-7cc8-4f60-94d0-174cc98310dc" WHERE "MUNICODE" = '%s'"""
    postgres_to_go = {
        "int4": "int",
        "text": "string",
        "float8": "float64",
        "date": "time.Time",
    }
    data = requests.get(url).json()
    raw_fields = data["result"]["fields"]
    fields = []
    for f in raw_fields:
        try:
            fields.append(
                to_field(
                    (
                        snake_case_to_pascal(f["id"]),
                        postgres_to_go[f["type"]],
                        'json:"%s"' % f["id"],
                    )
                )
            )
        # Debugging types when creating postgres_to_go
        except KeyError:
            if f["type"] not in ["tsvector", "geometry"]:
                raise KeyError(f["type"])

    return fields_to_struct("wprdcFields", fields)


def create_real_estate_portal_struct():
    # pages = ["GeneralInfo", "Building", "Tax", ""]
    pages = ["GeneralInfo"]
    urls = [
        "http://www2.alleghenycounty.us/RealEstate/%s.aspx?ParcelID=0945H00373000000"
        % page
        for page in pages
    ]
    responses = [requests.get(url) for url in urls]
    id_list = []
    for resp in responses:
        soup = bs4.BeautifulSoup(resp.text, "html.parser")
        attrs = [
            attr["id"]
            for attr in soup.findAll(id=lambda id: id and (id.find("lbl") != -1))
        ]
        id_list.extend(attrs)
    ids = set(id_list)
    fields = []
    for id in ids:
        if id.find("Text") != -1:
            continue
        fields.append(
            Field(
                name=id.split("lbl")[-1], type="scrapedHTML", annotation='id:"%s"' % id
            )
        )
    return fields_to_struct("realEstatePortal", fields)


def main():
    wprdc = create_wprdc_struct()
    real_estate_portal = create_real_estate_portal_struct()
    go = GoFile("external_data.go", "main")
    go.add_element(wprdc)
    go.add_element(real_estate_portal)
    go.write()

if __name__ == "__main__":
    main()
