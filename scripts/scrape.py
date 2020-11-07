import json
import re
import string
from typing import List, Dict

import requests
from bs4 import BeautifulSoup

TG_CORE_TYPES = ["String", "Boolean", "Integer", "Float"]

METHODS = "methods"
TYPES = "types"


def retrieve_api_info() -> Dict:
    r = requests.get("https://core.telegram.org/bots/api")
    soup = BeautifulSoup(r.text, features="html.parser")
    dev_rules = soup.find("div", {"id": "dev_page_content"})
    curr_type = ""
    curr_name = ""
    curr_desc = []

    items = {
        METHODS: dict(),
        TYPES: dict(),
    }

    for x in list(dev_rules.children):
        if x.name == "h3":
            # New category; clear name and type.
            curr_name = ""
            curr_type = ""
            curr_desc = []

        if x.name == "h4":
            name = x.find("a").get("name")
            if name and "-" in name:
                continue

            curr_name, curr_type = get_type_and_name(x, items)
            curr_desc = []

        if curr_type and curr_name and x.name == "p":
            description = x.get_text()
            # we only need returns for methods. If we have no no desription has been checked for returns yes
            if curr_type == METHODS and not curr_desc:
                get_method_return_type(curr_name, curr_type, description, items)

            curr_desc.append(description)
            items[curr_type][curr_name]["description"] = curr_desc

        if x.name == "table":
            get_fields(curr_name, curr_type, x, items)

    return items


def get_fields(curr_name, curr_type, x, items):
    body = x.find("tbody")
    fields = []
    for tr in body.find_all("tr"):
        children = list(tr.find_all("td"))
        if curr_type == TYPES and len(children) == 3:
            fields.append(
                {
                    "field": children[0].get_text(),
                    "types": clean_tg_type(children[1].get_text()),
                    "description": clean_tg_description(children[2].get_text()),

                }
            )

        elif curr_type == METHODS and len(children) == 4:
            fields.append(
                {
                    "parameter": children[0].get_text(),
                    "types": clean_tg_type(children[1].get_text()),
                    "required": children[2].get_text(),
                    "description": clean_tg_description(children[3].get_text()),
                }
            )

        else:
            print("idk what happened")
            print("Type", curr_type)
            print("Name", curr_name)
            print("n children", len(children))
            print(children)
            exit(1)
    items[curr_type][curr_name]["fields"] = fields


def get_method_return_type(curr_name, curr_type, description, items):
    ret_search = re.search("(?:on success,|returns)([^.]*)(?:on success)?", description, re.IGNORECASE)
    ret_search2 = re.search("([^.]*)(?:is returned)", description, re.IGNORECASE)
    if ret_search:
        extract_return_type(curr_type, curr_name, ret_search.group(1).strip(), items)
    elif ret_search2:
        extract_return_type(curr_type, curr_name, ret_search2.group(1).strip(), items)
    else:
        print("Failed to get return type for", curr_name)


def get_type_and_name(x, items):
    if x.text[0].isupper():
        curr_type = TYPES
    else:
        curr_type = METHODS
    curr_name = x.get_text()
    items[curr_type][curr_name] = {}

    return curr_name, curr_type


def extract_return_type(curr_type: str, curr_name: str, ret_str: str, items: Dict):
    array_match = re.search(r"(?:array of )+(\w*)", ret_str, re.IGNORECASE)
    if array_match:
        ret = clean_tg_type(array_match.group(1))
        rets = [f"Array of {r}" for r in ret]
        items[curr_type][curr_name]["returns"] = rets
    else:
        words = ret_str.split()
        rets = [
            r for ret in words
            for r in clean_tg_type(ret.translate(str.maketrans("", "", string.punctuation)))
            if ret[0].isupper()
        ]
        items[curr_type][curr_name]["returns"] = rets


def clean_tg_description(t: str) -> str:
    return t.replace('”', '"').replace('“', '"')


def get_proper_type(t: str) -> str:
    if t == "Messages":  # Avoids https://core.telegram.org/bots/api#sendmediagroup
        return "Message"

    elif t == "Float number":
        return "Float"

    elif t == "Int":
        return "Integer"

    elif t == "True" or t == "Bool":
        return "Boolean"

    return t


def clean_tg_type(t: str) -> List[str]:
    fixed_ors = [x.strip() for x in t.split(" or ")]  # Fix situations like "A or B"
    fixed_ands = [x.strip() for fo in fixed_ors for x in fo.split(" and ")]  # Fix situations like "A and B"
    fixed_commas = [x.strip() for fa in fixed_ands for x in fa.split(", ")]  # Fix situations like "A, B"
    return [get_proper_type(x) for x in fixed_commas]


def verify_type_parameters(items: Dict):
    for t, values in items[TYPES].items():
        # check all parameter types are valid
        for param in values.get("fields", []):
            types = param.get("types")
            for t in types:
                while t.startswith("Array of "):
                    t = t[len("Array of "):]

            if t not in items[TYPES] and t not in TG_CORE_TYPES:
                print("UNKNOWN FIELD TYPE", t)


def verify_method_parameters(items: Dict):
    # Type check all methods
    for method, values in items[METHODS].items():
        # check all methods have a return
        if not values.get("returns"):
            print(f"{method} has no return types!")
            continue

        if len(values.get("returns")) > 1:
            print(f"{method} has multiple return types: {values.get('returns')}")

        # check all parameter types are valid
        for param in values.get("fields", []):
            types = param.get("types")
            for t in types:
                while t.startswith("Array of "):
                    t = t[len("Array of "):]

                if t not in items[TYPES] and t not in TG_CORE_TYPES:
                    print("UNKNOWN PARAM TYPE", t)

        # check all return types are valid
        for ret in values.get("returns", []):
            while ret.startswith("Array of "):
                ret = ret[len("Array of "):]

            if ret not in items[TYPES] and ret not in TG_CORE_TYPES:
                print("UNKNOWN RETURN TYPE", ret)


if __name__ == '__main__':
    ITEMS = retrieve_api_info()
    verify_type_parameters(ITEMS)
    verify_method_parameters(ITEMS)

    with open("api.json", "w") as f:
        json.dump(ITEMS, f, indent=2)
