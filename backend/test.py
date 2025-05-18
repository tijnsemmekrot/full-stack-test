import os
from supabase import create_client, Client
import supabase
from fastapi import FastAPI

app = FastAPI()


def initDB():
    url: str = "https://ucexhnksccudzjcikyzi.supabase.co"
    key: str = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InVjZXhobmtzY2N1ZHpqY2lreXppIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDc0Nzg1NTQsImV4cCI6MjA2MzA1NDU1NH0.ys7zSEM3o0XkrOZbUUPA5k-EcatOrHxsIO36G3Z938M"
    supabase: Client = create_client(url, key)
    return supabase


def addNametoDB():
    supabase = initDB()
    supabase.table("names").insert({"name": "Julia-test"}).execute()


@app.get("/")
def main():
    supabase = initDB()
    addNametoDB()

    response = supabase.table("names").select("*").execute()
    print(response.data)


if __name__ == "__main__":
    main()
