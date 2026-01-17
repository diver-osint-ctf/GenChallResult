import { parse } from "@std/csv/parse";
import type { Chall, Challs } from "./types.ts";
import { CsvStream } from "https://deno.land/std@0.170.0/encoding/csv/stream.ts";

export async function loadChalls(filePath: string): Promise<Challs> {
  const file = await Deno.readTextFile(filePath);

  const parsedData = parse(file, {
    skipFirstRow: true,
  });

  const challs: Challs = {};

  for (const row of parsedData) {
    const id = row.id;
    const genre = row.category;
    const name = row.name;
    try {
      const score = Number(row.value);
      const solver = Number(0);
      const chall: Chall = {
        name,
        genre,
        score,
        solver,
      };
      challs[id] = chall;
    } catch (e) {
      console.log(`Error: ${e}`);
    }
  }
  return challs;
}

export async function loadTeams(filePath: string): Promise<Set<string>> {
  const file = await Deno.readTextFile(filePath);
  const parsedData = parse(file, {
    skipFirstRow: true,
  });

  const ignoredTeams = new Set<string>();

  for (const row of parsedData) {
    const id = row.id;
    const hidden = row.hidden;
    const banned = row.banned;

    if (hidden === "True" || banned === "True") {
      ignoredTeams.add(id);
    }
  }
  return ignoredTeams;
}

// solvesはファイルサイズが大きいため、ストリーミング処理を行う
export async function loadSolves(
  _filePath: string,
  challs: Challs,
  ignoredTeams: Set<string> = new Set(),
) {
  let filePath = _filePath;
  if (
    !filePath.startsWith("./") &&
    !filePath.startsWith("../") &&
    !filePath.startsWith("/")
  ) {
    filePath = `./${filePath}`;
  }
  const path = new URL(import.meta.resolve(filePath));
  const { readable } = await Deno.open(path);

  const data = readable
    .pipeThrough(new TextDecoderStream())
    .pipeThrough(new CsvStream());

  let passedHeader = false;
  for await (const line of data) {
    if (!passedHeader) {
      passedHeader = true;
      continue;
    }
    const challId = line[0];
    const teamId = line[2];

    if (ignoredTeams.has(teamId)) {
      continue;
    }

    if (challs[challId]) {
      challs[challId].solver++;
    }
  }
  return challs;
}
