"use client";

import { useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";

export default function BackButton({ text }) {
  const router = useRouter();

  return (
    <button
      type="button"
      onClick={() => router.back()}
      className="
        cursor-pointer flex items-center gap-2
        rounded-sm px-3 py-2 text-sm font-medium
        border border-[#60A5FA]/30
        bg-[#0F1016]
        text-[#BFDBFE]
        hover:bg-[#60A5FA]/10
        transition
      "
    >
      <ArrowLeft size={16} />
      {text}
    </button>
  );
}
