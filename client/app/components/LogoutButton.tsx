"use client";

import { webApi } from "@/app/lib/api/client";
import { LogOut } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function LogoutButton() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  async function onLogout() {
    setLoading(true);
    try {
      await webApi.post("/auth/logout");
      router.replace("/login");
      router.refresh();
    } catch {
    } finally {
      setLoading(false);
    }
  }

  return (
    <button
      type="button"
      onClick={onLogout}
      disabled={loading}
      className="
        flex items-center gap-2 cursor-pointer rounded-sm px-3 py-2 text-sm
        border border-white/10 bg-white/5
        text-[#A1A1AA]
        hover:text-white hover:bg-white/10
        transition
        disabled:opacity-50 disabled:cursor-not-allowed
      "
    >
      <LogOut size={16} />
      {loading ? "Logging out…" : "Logout"}
    </button>
  );
}
