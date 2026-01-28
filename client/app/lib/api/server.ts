import axios from "axios";
import { cookies } from "next/headers";

const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL;

export async function getServerApi() {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;

  return axios.create({
    baseURL,
    timeout: 15000,
    headers: token ? { Cookie: `access_token=${token}` } : undefined,
  });
}
