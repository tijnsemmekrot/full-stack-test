import os
from supabase import create_client, Client
from fastapi import FastAPI
from pydantic import BaseModel
from dotenv import load_dotenv
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()
load_dotenv()

origins = [
    "https://full-stack-test-chi.vercel.app",
    "http://localhost:3000",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


SUPABASE_URL = os.getenv("SUPABASE_URL")
SUPABASE_KEY = os.getenv("SUPABASE_KEY")

if not SUPABASE_URL or not SUPABASE_KEY:
    raise ValueError("Missing SUPABASE_URL or SUPABASE_KEY in environment variables")

supabase: Client = create_client(SUPABASE_URL, SUPABASE_KEY)


class NamePayLoad(BaseModel):
    name: str


@app.post("/api/firstName")
async def fetch_request(payload: NamePayLoad):
    print("Received name:", payload.name)


def addNametoDB():
    supabase.table("names").insert({"name": "Julia-test"}).execute()


@app.get("/")
def main():
    response = supabase.table("names").select("*").execute()
    print(response.data)


if __name__ == "__main__":
    main()
