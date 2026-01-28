"use client";

import { Eye, EyeOff } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { webApi } from "../lib/api/client";

export default function Login() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [form, setForm] = useState({ email: "", password: "" });
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);

    try {
      await webApi.post("/auth/login", form, {
        headers: { "Content-Type": "application/json" },
      });

      router.push("/competition");
      router.refresh();
    } catch {
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-[#0B0C12] flex items-center justify-center px-4">
      <div className="w-full max-w-sm rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-6 space-y-6">
        {/* Header */}
        <div className="text-center space-y-2">
          <h1 className="text-3xl sm:text-4xl font-semibold text-white">
            Log in
          </h1>
          <p className="text-sm text-[#A1A1AA]">
            Don&apos;t have a profile yet?{" "}
            <Link
              href="/register"
              className="text-[#60A5FA] underline underline-offset-4 hover:opacity-80 transition"
            >
              Create a profile
            </Link>
          </p>
        </div>

        {/* Form */}
        <form onSubmit={onSubmit} className="space-y-4">
          {/* Email */}
          <div className="space-y-1.5">
            <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
              Email
            </label>
            <input
              type="email"
              required
              value={form.email}
              onChange={(e) =>
                setForm((p) => ({ ...p, email: e.target.value }))
              }
              placeholder="Email"
              className="
                w-full rounded-xl px-3 py-2
                bg-[#0F1016]/80 border border-white/10
                text-white placeholder:text-white/30
                focus:outline-none focus:ring-2 focus:ring-[#60A5FA]/60
                transition
              "
            />
          </div>

          {/* Password */}
          <div className="space-y-1.5">
            <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
              Password
            </label>

            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                required
                value={form.password}
                onChange={(e) =>
                  setForm((p) => ({ ...p, password: e.target.value }))
                }
                placeholder="Password"
                className="
                  w-full rounded-xl px-3 py-2 pr-10
                  bg-[#0F1016]/80 border border-white/10
                  text-white placeholder:text-white/30
                  focus:outline-none focus:ring-2 focus:ring-[#A855F7]/60
                  transition
                "
              />

              <button
                type="button"
                aria-label={showPassword ? "Hide password" : "Show password"}
                onClick={() => setShowPassword((v) => !v)}
                className="
                  absolute inset-y-0 right-3 flex items-center
                  text-white/40 hover:text-white transition
                "
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>
          </div>

          {/* Submit */}
          <button
            type="submit"
            disabled={loading}
            className="
              w-full rounded-sm py-3 font-semibold
              bg-gradient-to-r from-[#A855F7] to-[#60A5FA]
              hover:opacity-90 transition
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            {loading ? "Logging in..." : "Log in"}
          </button>
        </form>
      </div>
    </div>
  );
}
