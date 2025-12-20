import { config } from "dotenv-flow";
import { NextResponse } from "next/server";
import { dirname, resolve } from "path";
import { Pool } from "pg";
import { fileURLToPath } from "url";

// Get the directory of this file
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Resolve the parent .env file from the dashboard directory
config({ path: resolve(__dirname, "../../../../../") });

const pool = new Pool({
  connectionString: process.env.NEXT_PUBLIC_DB_PATH,
});

export async function GET() {
  try {
    const now = new Date();
    const today = new Date(now.getTime() - now.getTimezoneOffset() * 60000)
      .toISOString()
      .slice(0, 10);

    const result = await pool.query(
      `WITH ranked AS (
         SELECT name, minutes_used, source,
                ROW_NUMBER() OVER (PARTITION BY source ORDER BY minutes_used DESC) as rn
         FROM daily_usage
         WHERE usage_date = $1
       )
       SELECT name, minutes_used, source
       FROM ranked
       WHERE rn <= 10
       ORDER BY source, minutes_used DESC`,
      [today]
    );
    return NextResponse.json(result.rows);
  } catch (error) {
    console.error("Database error:", error);
    return NextResponse.json(
      { error: "Failed to fetch usage data: " + (error as Error).message },
      { status: 500 }
    );
  }
}
