import axios from "axios";
import { toast } from "sonner";

const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL;

type ErrorResponse = {
  error?: string;
};

export const webApi = axios.create({
  baseURL,
  withCredentials: true,
  timeout: 15000,
});

webApi.interceptors.response.use(
  (res) => res,
  (err) => {
    if (!err?.response) {
      toast.error("Network error. Please try again.");
      return Promise.reject(err);
    }

    const data = err.response.data as ErrorResponse | undefined;
    const msg =
      typeof data?.error === "string" && data.error.trim()
        ? data.error.trim()
        : "Something went wrong.";

    toast.error(msg);
    return Promise.reject(err);
  },
);
