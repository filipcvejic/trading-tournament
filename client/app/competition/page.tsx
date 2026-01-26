import { redirect } from "next/navigation";
import { getServerApi } from "../lib/api/server";

export default async function CompetitionEntryPage() {
  const api = await getServerApi();
  const { data } = await api.get("/competitions/current");

  redirect(`/competition/${data.id}`);
}
