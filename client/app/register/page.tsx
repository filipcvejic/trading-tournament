"use client";

import { useMemo, useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { webApi } from "../lib/api/client";

/* ---------- FE validation helpers (UX only) ---------- */

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

function validateNoWhitespaceMinMax(value: string, min: number, max: number) {
  if (value.length < min) return `At least ${min} characters`;
  if (value.length > max) return `Max ${max} characters`;
  if (/\s/.test(value)) return "Must not contain whitespace";
  return null;
}

/* ---------- component ---------- */

export default function Register() {
  const router = useRouter();

  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const [form, setForm] = useState({
    username: "",
    discordUsername: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const usernameError = useMemo(
    () =>
      form.username ? validateNoWhitespaceMinMax(form.username, 3, 20) : null,
    [form.username],
  );

  const discordUsernameError = useMemo(
    () =>
      form.discordUsername
        ? validateNoWhitespaceMinMax(form.discordUsername, 2, 32)
        : null,
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
    emailValid;

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!canSubmit) return;

    setLoading(true);
    try {
      await webApi.post("/auth/register", {
        username: form.username,
        discordUsername: form.discordUsername,
        email: form.email,
        password: form.password,
      });

      router.push("/login");
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
            Create profile
          </h1>
          <p className="text-sm text-[#A1A1AA]">
            Already have an account?{" "}
            <Link
              href="/login"
              className="text-[#60A5FA] underline underline-offset-4 hover:opacity-80 transition"
            >
              Log in
            </Link>
          </p>
        </div>

        {/* Form */}
        <form onSubmit={onSubmit} className="space-y-4">
          {/* Username */}
          <div className="space-y-1.5">
            <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
              Username
            </label>
            <input
              value={form.username}
              onChange={(e) =>
                setForm((p) => ({ ...p, username: e.target.value }))
              }
              className="
                w-full rounded-xl px-3 py-2
                bg-[#0F1016]/80 border border-white/10
                text-white placeholder:text-white/30
                focus:outline-none focus:ring-2 focus:ring-[#60A5FA]/60
                transition
              "
              required
              placeholder="Username"
            />
            {usernameError && (
              <p className="text-xs text-red-400">{usernameError}</p>
            )}
          </div>

          {/* Discord username */}
          <div className="space-y-1.5">
            <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
              Discord username
            </label>
            <input
              value={form.discordUsername}
              onChange={(e) =>
                setForm((p) => ({ ...p, discordUsername: e.target.value }))
              }
              className="
                w-full rounded-xl px-3 py-2
                bg-[#0F1016]/80 border border-white/10
                text-white placeholder:text-white/30
                focus:outline-none focus:ring-2 focus:ring-[#A855F7]/60
                transition
              "
              required
              placeholder="Discord username"
            />
            {discordUsernameError && (
              <p className="text-xs text-red-400">{discordUsernameError}</p>
            )}
          </div>

          {/* Email */}
          <div className="space-y-1.5">
            <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
              Email
            </label>
            <input
              type="email"
              value={form.email}
              onChange={(e) =>
                setForm((p) => ({ ...p, email: e.target.value }))
              }
              className="
                w-full rounded-xl px-3 py-2
                bg-[#0F1016]/80 border border-white/10
                text-white placeholder:text-white/30
                focus:outline-none focus:ring-2 focus:ring-[#60A5FA]/60
                transition
              "
              required
              placeholder="Email"
            />
            {form.email.length > 0 && !emailValid && (
              <p className="text-xs text-red-400">
                Please enter a valid email address
              </p>
            )}
          </div>

          {/* Password */}
          <div className="space-y-1.5">
            <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
              Password
            </label>

            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                value={form.password}
                onChange={(e) =>
                  setForm((p) => ({ ...p, password: e.target.value }))
                }
                className="
                  w-full rounded-xl px-3 py-2 pr-10
                  bg-[#0F1016]/80 border border-white/10
                  text-white placeholder:text-white/30
                  focus:outline-none focus:ring-2 focus:ring-[#A855F7]/60
                  transition
                "
                required
                placeholder="Password"
              />

              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                className="absolute inset-y-0 right-3 flex items-center text-white/40 hover:text-white transition"
                aria-label={showPassword ? "Hide password" : "Show password"}
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>

            <ul className="mt-2 space-y-1 text-xs text-[#A1A1AA]">
              <li className={pwRules.length ? "text-[#BFDBFE]" : ""}>
                {pwRules.length ? "✓" : "•"} At least 11 characters
              </li>
              <li className={pwRules.lower ? "text-[#BFDBFE]" : ""}>
                {pwRules.lower ? "✓" : "•"} 1 lowercase letter
              </li>
              <li className={pwRules.upper ? "text-[#BFDBFE]" : ""}>
                {pwRules.upper ? "✓" : "•"} 1 uppercase letter
              </li>
              <li className={pwRules.special ? "text-[#BFDBFE]" : ""}>
                {pwRules.special ? "✓" : "•"} 1 special character
              </li>
            </ul>
          </div>

          {/* Confirm password */}
          <div className="space-y-1.5">
            <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
              Confirm password
            </label>
            <input
              type={showPassword ? "text" : "password"}
              value={form.confirmPassword}
              onChange={(e) =>
                setForm((p) => ({ ...p, confirmPassword: e.target.value }))
              }
              className="
                w-full rounded-xl px-3 py-2
                bg-[#0F1016]/80 border border-white/10
                text-white placeholder:text-white/30
                focus:outline-none focus:ring-2 focus:ring-[#60A5FA]/60
                transition
              "
              required
              placeholder="Confirm password"
            />
            {form.confirmPassword.length > 0 && !confirmValid && (
              <p className="text-xs text-red-400">Passwords do not match</p>
            )}
          </div>

          {/* Submit */}
          <button
            type="submit"
            disabled={!canSubmit}
            className="
              w-full rounded-sm py-3 font-semibold
              bg-gradient-to-r from-[#A855F7] to-[#60A5FA]
              hover:opacity-90 transition
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            {loading ? "Creating..." : "Create profile"}
          </button>
        </form>
      </div>
    </div>
  );
}
