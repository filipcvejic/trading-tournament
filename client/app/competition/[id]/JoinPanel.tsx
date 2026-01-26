"use client";

import { webApi } from "@/app/lib/api/client";
import { useMemo, useState } from "react";

type MeState = {
  hasRequestedAccount: boolean;
  hasJoined: boolean;
};

type FieldErrors = {
  login?: string;
  investorPassword?: string;
  broker?: string;
};

const NOTICE_REQUESTED =
  "You will receive an email soon with the trading account details required to join the competition.";

const NOTICE_JOIN =
  "Check your email and fill out this form with the trading account details you received.";

const NOTICE_JOINED =
  "You have successfully joined the competition. See you at kickoff! 🚀";

export default function JoinPanel({
  competitionId,
  initialMe,
}: {
  competitionId: string;
  initialMe: MeState;
}) {
  const [me, setMe] = useState<MeState>(initialMe);
  const [loading, setLoading] = useState(false);
  const [showJoin, setShowJoin] = useState(false);

  const [notice, setNotice] = useState<string | null>(
    initialMe.hasJoined
      ? NOTICE_JOINED
      : initialMe.hasRequestedAccount
        ? NOTICE_REQUESTED
        : null,
  );

  const [errors, setErrors] = useState<FieldErrors>({});

  const canRequest = useMemo(
    () => !me.hasRequestedAccount && !me.hasJoined,
    [me],
  );
  const canJoin = useMemo(() => me.hasRequestedAccount && !me.hasJoined, [me]);

  async function refetchMe() {
    const { data } = await webApi.get(`/competitions/${competitionId}/me`);
    setMe(data);
  }

  async function requestAccount() {
    setLoading(true);
    setErrors({});
    try {
      await webApi.post(`/competitions/${competitionId}/account-requests`);
      await refetchMe();
      setNotice(NOTICE_REQUESTED);
    } finally {
      setLoading(false);
    }
  }

  async function submitJoin(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setLoading(true);
    setErrors({});

    try {
      const form = new FormData(e.currentTarget);

      const loginStr = String(form.get("login") || "").trim();
      const investorPassword = String(
        form.get("investorPassword") || "",
      ).trim();
      const broker = String(form.get("broker") || "").trim();

      const nextErrors: FieldErrors = {};

      if (!loginStr) nextErrors.login = "Login is required.";
      else if (!/^\d+$/.test(loginStr))
        nextErrors.login = "Login must contain digits only.";

      if (!investorPassword)
        nextErrors.investorPassword = "Investor password is required.";
      if (!broker) nextErrors.broker = "Server is required.";

      if (Object.keys(nextErrors).length > 0) {
        setErrors(nextErrors);
        return;
      }

      const payload = {
        login: Number(loginStr),
        investorPassword,
        broker,
      };

      await webApi.post(`/competitions/${competitionId}/join`, payload);

      setMe((p) => ({ ...p, hasJoined: true }));
      setShowJoin(false);
      setNotice(NOTICE_JOINED);
      e.currentTarget.reset();
    } finally {
      setLoading(false);
    }
  }

  if (me.hasJoined) {
    return (
      <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-6 text-center">
        <p className="text-sm text-green-400">{NOTICE_JOINED}</p>
      </div>
    );
  }

  return (
    <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-6">
      <div className="text-center space-y-2">
        <h2 className="text-lg sm:text-xl font-semibold">Join setup</h2>
        <p className="text-sm text-[#A1A1AA]">
          Before the competition starts, request your account and then join
          using the credentials from email.
        </p>
      </div>

      {notice && (
        <div
          className="
            mt-4 mx-auto w-full max-w-xl rounded-xl
            border border-[#60A5FA]/30
            bg-[#0F1016]/80
            px-4 py-3
            text-sm text-[#BFDBFE]
            text-center
          "
        >
          {notice}
        </div>
      )}

      <div className="mt-6 flex flex-col items-center gap-4">
        {canRequest && (
          <button
            onClick={requestAccount}
            disabled={loading}
            className="
              w-full max-w-sm rounded-sm py-3 font-semibold
              bg-gradient-to-r from-[#A855F7] to-[#60A5FA]
              cursor-pointer hover:opacity-90 transition
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            {loading ? "Requesting..." : "Request an account"}
          </button>
        )}

        {canJoin && !showJoin && (
          <button
            onClick={() => {
              setErrors({});
              setShowJoin(true);
              setNotice(NOTICE_JOIN);
            }}
            disabled={loading}
            className="
              w-full max-w-sm rounded-sm py-3 font-semibold
              bg-gradient-to-r from-[#60A5FA] to-[#A855F7]
              hover:opacity-90 transition cursor-pointer
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            Join
          </button>
        )}

        {canJoin && showJoin && (
          <form
            onSubmit={submitJoin}
            aria-busy={loading}
            className={`w-full max-w-sm space-y-4 ${loading ? "opacity-90" : ""}`}
          >
            <div>
              <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
                Login
              </label>
              <input
                name="login"
                type="text"
                inputMode="numeric"
                autoComplete="off"
                placeholder="e.g. 12345678"
                required
                disabled={loading}
                onChange={(e) => {
                  const digitsOnly = e.target.value.replace(/\D/g, "");
                  if (digitsOnly !== e.target.value)
                    e.target.value = digitsOnly;

                  setErrors((prev) => ({
                    ...prev,
                    login: digitsOnly ? undefined : prev.login,
                  }));
                }}
                className="
                  w-full rounded-xl px-3 py-2
                  bg-[#0F1016] border border-white/10
                  focus:outline-none focus:ring-2 focus:ring-[#A855F7]/60 focus:border-[#A855F7]/50
                "
              />
              {errors.login ? (
                <p className="mt-1 text-xs text-red-400">{errors.login}</p>
              ) : null}
            </div>

            <div>
              <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
                Investor password
              </label>
              <input
                name="investorPassword"
                placeholder="Investor password"
                required
                disabled={loading}
                type="password"
                onChange={(e) =>
                  setErrors((prev) => ({
                    ...prev,
                    investorPassword: e.target.value.trim()
                      ? undefined
                      : prev.investorPassword,
                  }))
                }
                className="
                  w-full rounded-xl px-3 py-2
                  bg-[#0F1016] border border-white/10
                  focus:outline-none focus:ring-2 focus:ring-[#A855F7]/60 focus:border-[#A855F7]/50
                "
              />
              {errors.investorPassword ? (
                <p className="mt-1 text-xs text-red-400">
                  {errors.investorPassword}
                </p>
              ) : null}
            </div>

            <div>
              <label className="block mb-1.5 text-sm font-medium text-[#C7D2FE]">
                Server
              </label>
              <input
                name="broker"
                placeholder="Broker"
                required
                disabled={loading}
                onChange={(e) =>
                  setErrors((prev) => ({
                    ...prev,
                    broker: e.target.value.trim() ? undefined : prev.broker,
                  }))
                }
                className="
                  w-full rounded-xl px-3 py-2
                  bg-[#0F1016] border border-white/10
                  focus:outline-none focus:ring-2 focus:ring-[#A855F7]/60 focus:border-[#A855F7]/50
                "
              />
              {errors.broker ? (
                <p className="mt-1 text-xs text-red-400">{errors.broker}</p>
              ) : null}
            </div>

            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => {
                  setShowJoin(false);
                  setErrors({});
                  setNotice(me.hasRequestedAccount ? NOTICE_REQUESTED : null);
                }}
                disabled={loading}
                className="
                  flex-1 rounded-sm py-2.5 font-semibold
                  border border-white/10 bg-black/20
                  hover:bg-white/5 transition
                  disabled:opacity-50 disabled:cursor-not-allowed
                "
              >
                Cancel
              </button>

              <button
                type="submit"
                disabled={loading}
                className="
                  flex-1 rounded-sm py-2.5 font-semibold
                  bg-gradient-to-r from-[#A855F7] to-[#60A5FA]
                  hover:opacity-90 transition cursor-pointer
                  disabled:opacity-50 disabled:cursor-not-allowed
                "
              >
                {loading ? "Joining..." : "Submit"}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
