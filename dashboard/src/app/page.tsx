"use client";

import UsageChart from "@/components/UsageChart";
import { useEffect, useState } from "react";

interface Usage {
  name: string;
  minutes_used: number;
  source: string;
}

type ChartType = "bar" | "line" | "pie";

export default function Home() {
  const [usageData, setUsageData] = useState<Usage[]>([]);
  const [loading, setLoading] = useState(true);
  const [chartType, setChartType] = useState<ChartType>("bar");

  useEffect(() => {
    fetch("/api/usage")
      .then((res) => res.json())
      .then((data) => setUsageData(data))
      .catch((err) => console.error("Failed to fetch usage data:", err))
      .finally(() => setLoading(false));
  }, []);

  const chartTypes: { value: ChartType; label: string }[] = [
    { value: "bar", label: "Bar" },
    { value: "line", label: "Line" },
    { value: "pie", label: "Pie" },
  ];

  return (
    <main className="min-h-screen bg-gray-900 p-6 sm:p-8 lg:p-12">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="text-center mb-8 lg:mb-12">
          <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-3 text-white">
            Focus Tracker
          </h1>
        </div>

        {/* Chart Type Selector */}
        <div className="flex justify-center mb-8">
          <div className="inline-flex items-center gap-2 p-1.5 bg-gray-800 rounded-2xl shadow-lg ring-1 ring-gray-700">
            {chartTypes.map((type) => (
              <button
                key={type.value}
                onClick={() => setChartType(type.value)}
                className={`
                  px-4 py-2.5 rounded-xl font-medium text-sm transition-all duration-200
                  flex items-center gap-2
                  ${
                    chartType === type.value
                      ? "bg-blue-600 text-white shadow-md"
                      : "text-gray-400 hover:bg-gray-700 hover:text-white"
                  }
                `}
              >
                <span>{type.label}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Content */}
        {loading ? (
          <div className="flex flex-col items-center justify-center py-20">
            <div className="animate-spin rounded-full h-12 w-12 border-4 border-gray-700 border-t-blue-500 mb-4"></div>
            <p className="text-gray-400 text-sm">Loading your data...</p>
          </div>
        ) : usageData.length === 0 ? (
          <div className="text-center py-20">
            <div className="text-6xl mb-4">📊</div>
            <p className="text-gray-400 text-lg">No data found for today.</p>
            <p className="text-gray-500 text-sm mt-2">
              Start tracking to see your productivity insights!
            </p>
          </div>
        ) : (
          <UsageChart data={usageData} chartType={chartType} />
        )}
      </div>
    </main>
  );
}
