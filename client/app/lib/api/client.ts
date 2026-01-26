import axios from "axios";

const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL;

export const webApi = axios.create({
  baseURL,
  withCredentials: true,
});
