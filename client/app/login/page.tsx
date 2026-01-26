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
  const [error, setError] = useState("");

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await webApi.post("/auth/login", form, {
        withCredentials: true,
        headers: {
          "Content-Type": "application/json",
        },
      });

      router.push("/competition");
    } catch (err: any) {
      setError("Invalid email or password");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="bg-[#1e1e1e] min-h-screen flex items-center justify-center">
      <div className="w-full max-w-sm space-y-6">
        {/* Header */}
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-semibold">Log in</h1>
          <p className="text-sm text-gray-300">
            Don&apos;t have a profile yet?{" "}
            <Link
              href="/register"
              className="text-[#83c0ff] underline underline-offset-4 hover:opacity-80 transition"
            >
              Create a profile
            </Link>
          </p>
        </div>

        {/* Form */}
        <form onSubmit={onSubmit} className="space-y-4">
          {/* Email */}
          <div className="space-y-1.5">
            <label htmlFor="email" className="block font-medium">
              Email
            </label>
            <input
              id="email"
              type="email"
              placeholder="Email"
              required
              value={form.email}
              onChange={(e) =>
                setForm((p) => ({ ...p, email: e.target.value }))
              }
              className="
            w-full rounded-xl px-3 py-2
            bg-[#36373b] border-2 border-transparent
            focus:border-[#83c0ff] focus:outline-none
            transition
          "
            />
          </div>

          {/* Password */}
          <div className="space-y-1.5">
            <label htmlFor="password" className="block font-medium">
              Password
            </label>

            <div className="relative">
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                placeholder="Password"
                required
                value={form.password}
                onChange={(e) =>
                  setForm((p) => ({ ...p, password: e.target.value }))
                }
                className="
              w-full rounded-xl px-3 py-2 pr-10
              bg-[#36373b] border-2 border-transparent
              focus:border-[#83c0ff] focus:outline-none
              transition
            "
              />

              <button
                type="button"
                aria-label={showPassword ? "Hide password" : "Show password"}
                onClick={() => setShowPassword((prev) => !prev)}
                className="
              absolute inset-y-0 right-3
              flex items-center
              text-gray-400 hover:text-white
              transition
            "
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>
          </div>

          {/* Error */}
          {error ? <p className="text-sm text-red-400">{error}</p> : null}

          {/* Forgot password */}
          <div className="flex justify-end">
            <a
              href="#"
              className="text-sm underline underline-offset-4 hover:opacity-80 transition"
            >
              Forgot password?
            </a>
          </div>

          {/* Submit */}
          <button
            type="submit"
            disabled={loading}
            className="
          w-full rounded-sm py-2.5
          bg-[#0781fe] font-medium
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
