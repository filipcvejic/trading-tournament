import Image from "next/image";
import JoinPanel from "./JoinPanel";
import LeaderboardTable from "./LeaderboardTable";
import { getServerApi } from "@/app/lib/api/server";
import { notFound } from "next/navigation";
import LogoutButton from "@/app/components/LogoutButton";

function isStarted(startsAt: string) {
  return Date.now() >= new Date(startsAt).getTime();
}
function isEnded(endsAt: string) {
  return Date.now() > new Date(endsAt).getTime();
}

export default async function CompetitionPage({
  params,
}: {
  params: { id: string };
}) {
  const { id } = await params;

  if (!id || id.length < 10) notFound();

  const api = await getServerApi();

  const [{ data: competition }, { data: me }] = await Promise.all([
    api.get(`/competitions/${id}`),
    api.get(`/competitions/${id}/me`),
  ]);

  const started = isStarted(competition.startsAt);
  const ended = isEnded(competition.endsAt);

  return (
    <div className="min-h-screen bg-[#0B0C10] text-white">
      {/* subtle neon glow */}
      <div className="absolute top-4 right-4">
        <LogoutButton />
      </div>
      <div className="pointer-events-none fixed inset-0 opacity-60">
        <div className="absolute -top-24 left-1/2 h-72 w-[42rem] -translate-x-1/2 rounded-full bg-[#A855F7]/20 blur-3xl" />
        <div className="absolute top-48 left-1/2 h-72 w-[42rem] -translate-x-1/2 rounded-full bg-[#60A5FA]/15 blur-3xl" />
      </div>

      <div className="relative mx-auto w-full max-w-6xl px-4 py-10 space-y-8">
        {/* Header */}
        <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur px-6 py-6">
          <div className="flex items-start gap-4">
            <div className="relative h-12 w-12 shrink-0 overflow-hidden rounded-xl border border-white/10 bg-black/30">
              <Image
                src="/logo.png"
                alt="Logo"
                fill
                className="object-cover"
                priority
              />
            </div>

            <div className="min-w-0">
              <h1 className="text-2xl sm:text-3xl font-semibold tracking-tight">
                {competition.name}
              </h1>
              <p className="mt-2 text-sm sm:text-base text-[#A1A1AA]">
                {competition.description}
              </p>

              <div className="mt-4 flex flex-wrap gap-2 text-xs sm:text-sm">
                <span className="rounded-full border border-white/10 bg-black/20 px-3 py-1 text-[#C7D2FE]">
                  Start:{" "}
                  <span className="text-white">
                    {new Date(competition.startsAt).toLocaleString()}
                  </span>
                </span>
                <span className="rounded-full border border-white/10 bg-black/20 px-3 py-1 text-[#C7D2FE]">
                  End:{" "}
                  <span className="text-white">
                    {new Date(competition.endsAt).toLocaleString()}
                  </span>
                </span>

                <span
                  className={[
                    "rounded-full border px-3 py-1",
                    ended
                      ? "border-white/10 bg-white/5 text-[#A1A1AA]"
                      : started
                        ? "border-[#60A5FA]/30 bg-[#60A5FA]/10 text-[#BFDBFE]"
                        : "border-[#A855F7]/30 bg-[#A855F7]/10 text-[#E9D5FF]",
                  ].join(" ")}
                >
                  {ended ? "Finished" : started ? "Live" : "Upcoming"}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* CTA (samo pre starta) */}
        {!started && (
          <JoinPanel competitionId={competition.id} initialMe={me} />
        )}

        {/* Leaderboard (od starta pa nadalje) */}
        {started && (
          <LeaderboardTable
            competitionId={competition.id}
            endsAt={competition.endsAt}
          />
        )}
      </div>
    </div>
  );
}
