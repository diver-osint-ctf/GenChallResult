import { loadChalls, loadSolves, loadTeams } from "./loadFiles.ts";
import { assertEquals } from "jsr:@std/assert";

Deno.test("loadChalls - should parse challenges CSV", async () => {
  const tmpFile = await Deno.makeTempFile({ suffix: ".csv" });
  await Deno.writeTextFile(
    tmpFile,
    "id,name,description,max_attempts,value,category,type,state,requirements\n" +
      "1,Test Challenge,,0,100,Web,standard,visible,\n" +
      "2,Another Challenge,,0,200,Crypto,standard,visible,\n",
  );

  try {
    const challs = await loadChalls(tmpFile);
    assertEquals(challs["1"].name, "Test Challenge");
    assertEquals(challs["1"].genre, "Web");
    assertEquals(challs["1"].score, 100);
    assertEquals(challs["1"].solver, 0);
    assertEquals(challs["2"].name, "Another Challenge");
    assertEquals(challs["2"].genre, "Crypto");
    assertEquals(challs["2"].score, 200);
    assertEquals(challs["2"].solver, 0);
  } finally {
    await Deno.remove(tmpFile);
  }
});

Deno.test("loadTeams - should identify hidden and banned teams", async () => {
  const tmpFile = await Deno.makeTempFile({ suffix: ".csv" });
  await Deno.writeTextFile(
    tmpFile,
    "id,oauth_id,name,email,password,secret,website,affiliation,country,bracket_id,hidden,banned,captain_id,created\n" +
      "1,,admin,,hash,,,,,,True,False,1,2026-01-01\n" +
      "2,,team1,,hash,,,,,,False,False,2,2026-01-01\n" +
      "3,,banned_team,,hash,,,,,,False,True,3,2026-01-01\n",
  );

  try {
    const ignoredTeams = await loadTeams(tmpFile);
    assertEquals(ignoredTeams.has("1"), true); // hidden
    assertEquals(ignoredTeams.has("2"), false); // normal
    assertEquals(ignoredTeams.has("3"), true); // banned
  } finally {
    await Deno.remove(tmpFile);
  }
});

Deno.test(
  "loadTeams - should return empty set when no hidden/banned teams",
  async () => {
    const tmpFile = await Deno.makeTempFile({ suffix: ".csv" });
    await Deno.writeTextFile(
      tmpFile,
      "id,oauth_id,name,email,password,secret,website,affiliation,country,bracket_id,hidden,banned,captain_id,created\n" +
        "1,,team1,,hash,,,,,,False,False,1,2026-01-01\n" +
        "2,,team2,,hash,,,,,,False,False,2,2026-01-01\n",
    );

    try {
      const ignoredTeams = await loadTeams(tmpFile);
      assertEquals(ignoredTeams.size, 0);
    } finally {
      await Deno.remove(tmpFile);
    }
  },
);

Deno.test("loadSolves - should count solves per challenge", async () => {
  const tmpFile = await Deno.makeTempFile({ suffix: ".csv" });
  await Deno.writeTextFile(
    tmpFile,
    "challenge_id,user_id,team_id,id,ip,provided,type,date\n" +
      "1,1,1,1,127.0.0.1,flag1,correct,2026-01-01\n" +
      "1,2,2,2,127.0.0.1,flag1,correct,2026-01-01\n" +
      "2,1,1,3,127.0.0.1,flag2,correct,2026-01-01\n",
  );

  const challs = {
    "1": { name: "C1", genre: "Web", score: 100, solver: 0 },
    "2": { name: "C2", genre: "Crypto", score: 200, solver: 0 },
  };

  try {
    const result = await loadSolves(tmpFile, challs);
    assertEquals(result["1"].solver, 2);
    assertEquals(result["2"].solver, 1);
  } finally {
    await Deno.remove(tmpFile);
  }
});

Deno.test("loadSolves - should exclude ignored teams", async () => {
  const tmpFile = await Deno.makeTempFile({ suffix: ".csv" });
  await Deno.writeTextFile(
    tmpFile,
    "challenge_id,user_id,team_id,id,ip,provided,type,date\n" +
      "1,1,1,1,127.0.0.1,flag1,correct,2026-01-01\n" +
      "1,2,2,2,127.0.0.1,flag1,correct,2026-01-01\n" +
      "1,3,3,3,127.0.0.1,flag1,correct,2026-01-01\n",
  );

  const challs = {
    "1": { name: "C1", genre: "Web", score: 100, solver: 0 },
  };
  const ignoredTeams = new Set(["1"]);

  try {
    const result = await loadSolves(tmpFile, challs, ignoredTeams);
    assertEquals(result["1"].solver, 2);
  } finally {
    await Deno.remove(tmpFile);
  }
});

Deno.test(
  "loadSolves - should skip solves for unknown challenges",
  async () => {
    const tmpFile = await Deno.makeTempFile({ suffix: ".csv" });
    await Deno.writeTextFile(
      tmpFile,
      "challenge_id,user_id,team_id,id,ip,provided,type,date\n" +
        "1,1,1,1,127.0.0.1,flag1,correct,2026-01-01\n" +
        "999,2,2,2,127.0.0.1,flag999,correct,2026-01-01\n",
    );

    const challs = {
      "1": { name: "C1", genre: "Web", score: 100, solver: 0 },
    };

    try {
      const result = await loadSolves(tmpFile, challs);
      assertEquals(result["1"].solver, 1);
      assertEquals(result["999"], undefined);
    } finally {
      await Deno.remove(tmpFile);
    }
  },
);
