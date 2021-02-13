import requests
from golang_code_generator.golang_code_gen import fields_to_struct, to_field


url = """https://data.wprdc.org/api/3/action/datastore_search_sql?sql=SELECT * FROM "518b583f-7cc8-4f60-94d0-174cc98310dc" WHERE "MUNICODE" = '%s'"""
postgres_to_go = {
    "int4": "int",
    "text": "string",
    "float8": "float64",
    "date": "time.Time",

}


def main():
    data = requests.get(url).json()
    raw_fields = data["result"]["fields"]
    fields = []
    for f in raw_fields:
        try:
            fields.append(to_field((f["id"], postgres_to_go[f["type"]])))
        except KeyError:
            if f["type"] not in ["tsvector", "geometry"]:
                raise KeyError(f["type"])

    print(fields_to_struct("wprdcFields", fields))


if __name__ == '__main__':
    main()
