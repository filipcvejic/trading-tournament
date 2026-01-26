"use client";

import { useMemo, useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { webApi } from "../lib/api/client";

/* ---------- validation helpers ---------- */

function validatePassword(pw: string) {
  return {
    length: pw.length >= 11,
    lower: /[a-z]/.test(pw),
    upper: /[A-Z]/.test(pw),
    special: /[^A-Za-z0-9]/.test(pw),
  };
}

function validateEmail(email: string) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

function validateUsername(username: string) {
  if (username.length < 3) return "At least 3 characters";
  if (username.length > 20) return "Max 20 characters";
  if (!/^[a-zA-Z0-9_]+$/.test(username))
    return "Only letters, numbers and underscore";
  return null;
}

/* ---------- component ---------- */

export default function Register() {
  const router = useRouter();

  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState({
    username: "",
    discordUsername: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const usernameError = useMemo(
    () => (form.username ? validateUsername(form.username) : null),
    [form.username],
  );

  const discordUsernameError = useMemo(
    () =>
      form.discordUsername ? validateUsername(form.discordUsername) : null,
    [form.discordUsername],
  );

  const pwRules = useMemo(
    () => validatePassword(form.password),
    [form.password],
  );

  const passwordValid =
    pwRules.length && pwRules.lower && pwRules.upper && pwRules.special;

  const confirmValid =
    form.confirmPassword.length > 0 && form.confirmPassword === form.password;

  const emailValid = useMemo(
    () => (form.email ? validateEmail(form.email) : false),
    [form.email],
  );

  const canSubmit =
    !loading &&
    !usernameError &&
    !discordUsernameError &&
    passwordValid &&
    confirmValid &&
    form.email.length > 0;

  /* ---------- submit ---------- */

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!canSubmit) return;

    setLoading(true);
    setError("");

    try {
      await webApi.post("/auth/register", {
        username: form.username,
        discordUsername: form.discordUsername,
        email: form.email,
        password: form.password,
      });

      router.push("/login");
    } catch {
      setError("Registration failed. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  /* ---------- UI ---------- */

  return (
    <div className="bg-[#1e1e1e] min-h-screen flex items-center justify-center">
      <div className="w-full max-w-sm space-y-6">
        {/* Header */}
        <div className="space-y-2 text-center">
          <h1 className="text-5xl font-semibold">Create profile</h1>
          <p className="text-sm text-gray-300">
            Already have an account?{" "}
            <Link
              href="/login"
              className="text-[#83c0ff] underline underline-offset-4 hover:opacity-80 transition"
            >
              Log in
            </Link>
          </p>
        </div>

        {/* Form */}
        <form onSubmit={onSubmit} className="space-y-4">
          {/* Username */}
          <div className="space-y-1.5">
            <label className="block font-medium">Username</label>
            <input
              value={form.username}
              onChange={(e) =>
                setForm((p) => ({ ...p, username: e.target.value }))
              }
              className="
                w-full rounded-xl px-3 py-2
                bg-[#36373b] border-2 border-transparent
                focus:border-[#83c0ff] focus:outline-none transition
              "
              required
            />
            {usernameError && (
              <p className="text-xs text-red-400">{usernameError}</p>
            )}
          </div>

          <div className="space-y-1.5">
            <label className="block font-medium">Discord username</label>
            <input
              value={form.discordUsername}
              onChange={(e) =>
                setForm((p) => ({ ...p, discordUsername: e.target.value }))
              }
              className="
                w-full rounded-xl px-3 py-2
                bg-[#36373b] border-2 border-transparent
                focus:border-[#83c0ff] focus:outline-none transition
              "
              required
            />
            {discordUsernameError && (
              <p className="text-xs text-red-400">{discordUsernameError}</p>
            )}
          </div>

          {/* Email */}
          <div className="space-y-1.5">
            <label className="block font-medium">Email</label>
            <input
              type="email"
              value={form.email}
              onChange={(e) =>
                setForm((p) => ({ ...p, email: e.target.value }))
              }
              className="
      w-full rounded-xl px-3 py-2
      bg-[#36373b] border-2 border-transparent
      focus:border-[#83c0ff] focus:outline-none transition
    "
              required
            />

            {form.email.length > 0 && !emailValid && (
              <p className="text-xs text-red-400">
                Please enter a valid email address
              </p>
            )}
          </div>

          {/* Password */}
          <div className="space-y-1.5">
            <label className="block font-medium">Password</label>

            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                value={form.password}
                onChange={(e) =>
                  setForm((p) => ({ ...p, password: e.target.value }))
                }
                className="
                  w-full rounded-xl px-3 py-2 pr-10
                  bg-[#36373b] border-2 border-transparent
                  focus:border-[#83c0ff] focus:outline-none transition
                "
                required
              />

              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                className="
                  absolute inset-y-0 right-3 flex items-center
                  text-gray-400 hover:text-white transition
                "
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>

            {/* Password rules */}
            <ul className="mt-2 space-y-1 text-xs text-gray-300">
              <li className={pwRules.length ? "text-green-400" : ""}>
                {pwRules.length ? "✓" : "•"} At least 11 characters
              </li>
              <li className={pwRules.lower ? "text-green-400" : ""}>
                {pwRules.lower ? "✓" : "•"} 1 lowercase letter
              </li>
              <li className={pwRules.upper ? "text-green-400" : ""}>
                {pwRules.upper ? "✓" : "•"} 1 uppercase letter
              </li>
              <li className={pwRules.special ? "text-green-400" : ""}>
                {pwRules.special ? "✓" : "•"} 1 special character
              </li>
            </ul>
          </div>

          {/* Confirm password */}
          <div className="space-y-1.5">
            <label className="block font-medium">Confirm password</label>
            <input
              type={showPassword ? "text" : "password"}
              value={form.confirmPassword}
              onChange={(e) =>
                setForm((p) => ({ ...p, confirmPassword: e.target.value }))
              }
              className="
                w-full rounded-xl px-3 py-2
                bg-[#36373b] border-2 border-transparent
                focus:border-[#83c0ff] focus:outline-none transition
              "
              required
            />
            {form.confirmPassword.length > 0 && !confirmValid && (
              <p className="text-xs text-red-400">Passwords do not match</p>
            )}
          </div>

          {/* Error */}
          {error && <p className="text-sm text-red-400">{error}</p>}

          {/* Submit */}
          <button
            type="submit"
            disabled={!canSubmit}
            className="
              w-full rounded-sm py-2.5
              bg-[#0781fe] font-medium
              hover:opacity-90 transition
              disabled:opacity-60 disabled:cursor-not-allowed
            "
          >
            {loading ? "Creating..." : "Create profile"}
          </button>
        </form>
      </div>
    </div>
  );
}
