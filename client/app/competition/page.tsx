import axios from "axios";
import { redirect } from "next/navigation";
import { getServerApi } from "../lib/api/server";
import LogoutButton from "../components/LogoutButton";

export default async function CompetitionEntryPage() {
  const api = await getServerApi();

  try {
    const { data } = await api.get("/competitions/current");
    redirect(`/competition/${data.id}`);
  } catch (err: any) {
    // Ako je redirect error, pusti ga da prođe (NE LOGUJ)
    if (err?.digest?.startsWith?.("NEXT_REDIRECT")) {
      throw err;
    }

    // Ako je Axios error i 404 => fallback UI, bez loga
    if (axios.isAxiosError(err) && err.response?.status === 404) {
      // show fallback UI below
    } else {
      // Ostalo: loguj
      console.error("Failed to load current competition:", err);
    }
  }

  return (
    <div className="min-h-screen bg-[#0B0C12] relative px-4">
      {/* Logout – top right */}
      <div className="absolute top-4 right-4 z-50">
        <LogoutButton />
      </div>

      {/* Page content */}
      <div className="flex justify-center pt-24">
        <div className="w-full max-w-md rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-6">
          <div className="text-center space-y-2">
            <h1 className="text-2xl sm:text-3xl font-semibold text-white">
              No active competition
            </h1>
            <p className="text-sm text-[#A1A1AA]">
              There is no active or upcoming competition right now. Please check
              back later.
            </p>
          </div>

          <div
            className="
            mt-5 mx-auto w-full max-w-xl rounded-xl
            border border-[#60A5FA]/30
            bg-[#0F1016]/80
            px-4 py-3
            text-sm text-[#BFDBFE]
            text-center
          "
          >
            We’ll announce the next one soon.
          </div>

          <div className="mt-6 flex justify-center">
            <span
              className="
              inline-flex items-center justify-center
              rounded-sm px-4 py-2 font-semibold
              bg-gradient-to-r from-[#A855F7] to-[#60A5FA]
              opacity-90
            "
            >
              Stay tuned
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
