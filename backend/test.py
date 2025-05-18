import os
from supabase import create_client, Client
from fastapi import FastAPI
from dotenv import load_dotenv

app = FastAPI()
load_dotenv()

SUPABASE_URL = os.getenv("SUPABASE_URL")
SUPABASE_KEY = os.getenv("SUPABASE_KEY")
print(SUPABASE_URL, SUPABASE_KEY)
supabase = create_client(SUPABASE_URL, SUPABASE_KEY)


def addNametoDB():
    supabase.table("names").insert({"name": "Julia-test"}).execute()


@app.get("/")
def main():
    addNametoDB()

    response = supabase.table("names").select("*").execute()
    print(response.data)


if __name__ == "__main__":
    main()
