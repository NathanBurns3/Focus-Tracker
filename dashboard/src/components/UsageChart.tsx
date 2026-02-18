import {
  Bar,
  BarChart,
  Line,
  LineChart,
  Pie,
  PieChart,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  CartesianGrid,
} from "recharts";

interface Usage {
  name: string;
  minutes_used: number;
  source: string;
}

type ChartType = "bar" | "line" | "pie";

export default function UsageChart({
  data,
  chartType = "bar",
}: {
  data: Usage[];
  chartType?: ChartType;
}) {
  const normalize = (items: Usage[]) =>
    items.map((d) => ({
      ...d,
      minutes_used: Number(d.minutes_used),
    }));

  const apps = normalize(data.filter((d) => d.source === "desktop"));
  const sites = normalize(data.filter((d) => d.source === "chrome"));

  // Color schemes
  const COLORS = {
    apps: ["#4F46E5", "#6366F1", "#818CF8", "#A5B4FC", "#C7D2FE"],
    sites: ["#10B981", "#34D399", "#6EE7B7", "#A7F3D0", "#D1FAE5"],
  };

  // Custom tooltip component for a cleaner, modern look
  const TooltipContent = (props: {
    active?: boolean;
    payload?: Array<{ value: number }>;
    label?: string;
  }) => {
    const { active, payload, label } = props;
    if (active && payload && payload.length) {
      const minutes = payload[0].value as number;
      return (
        <div className="rounded-2xl bg-gradient-to-br from-gray-900/95 to-gray-800/95 text-white px-4 py-3 shadow-2xl backdrop-blur-sm border border-white/20">
          <div className="text-xs text-gray-300 font-medium">{label}</div>
          <div className="mt-1 text-base font-bold">{minutes} min</div>
        </div>
      );
    }
    return null;
  };

  // Render chart based on type
  const renderChart = (items: Usage[], colors: string[], title: string) => {
    const gradientId = title.toLowerCase();

    switch (chartType) {
      case "bar":
        return (
          <BarChart
            data={items}
            margin={{ top: 8, right: 12, left: 0, bottom: 60 }}
          >
            <defs>
              <linearGradient
                id={`bar${gradientId}`}
                x1="0"
                y1="0"
                x2="0"
                y2="1"
              >
                <stop offset="0%" stopColor={colors[0]} stopOpacity={0.95} />
                <stop offset="100%" stopColor={colors[1]} stopOpacity={0.75} />
              </linearGradient>
            </defs>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="#374151"
              opacity={0.3}
            />
            <XAxis
              dataKey="name"
              tick={{ fill: "#9CA3AF", fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              interval={0}
              angle={-45}
              textAnchor="end"
              height={80}
            />
            <YAxis
              tick={{ fill: "#9CA3AF", fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              width={40}
            />
            <Tooltip
              content={<TooltipContent />}
              cursor={{ fill: "rgba(255,255,255,0.05)" }}
            />
            <Bar
              dataKey="minutes_used"
              fill={`url(#bar${gradientId})`}
              radius={[10, 10, 0, 0]}
              barSize={32}
            />
          </BarChart>
        );

      case "line":
        return (
          <LineChart
            data={items}
            margin={{ top: 8, right: 12, left: 0, bottom: 60 }}
          >
            <defs>
              <linearGradient
                id={`line${gradientId}`}
                x1="0"
                y1="0"
                x2="1"
                y2="0"
              >
                <stop offset="0%" stopColor={colors[0]} stopOpacity={1} />
                <stop offset="100%" stopColor={colors[1]} stopOpacity={1} />
              </linearGradient>
            </defs>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="#374151"
              opacity={0.3}
            />
            <XAxis
              dataKey="name"
              tick={{ fill: "#9CA3AF", fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              interval={0}
              angle={-45}
              textAnchor="end"
              height={80}
            />
            <YAxis
              tick={{ fill: "#9CA3AF", fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              width={40}
            />
            <Tooltip content={<TooltipContent />} />
            <Line
              type="monotone"
              dataKey="minutes_used"
              stroke={`url(#line${gradientId})`}
              strokeWidth={3}
              dot={{ fill: colors[0], strokeWidth: 2, r: 5 }}
              activeDot={{ r: 7, strokeWidth: 2 }}
            />
          </LineChart>
        );

      case "pie":
        return (
          <PieChart margin={{ top: 8, right: 12, left: 0, bottom: 0 }}>
            <Pie
              data={items}
              cx="50%"
              cy="50%"
              labelLine={false}
              label={(entry) => `${entry.name}`}
              outerRadius={130}
              dataKey="minutes_used"
            >
              {items.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={colors[index % colors.length]}
                />
              ))}
            </Pie>
            <Tooltip content={<TooltipContent />} />
          </PieChart>
        );

      default:
        return null;
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 lg:gap-8">
      {/* Apps */}
      <div className="group p-6 sm:p-8 bg-gray-800 rounded-3xl shadow-xl ring-1 ring-gray-700 hover:shadow-2xl transition-all duration-300">
        <div className="flex items-center gap-3 mb-6">
          <div className="p-2.5 bg-blue-900/30 rounded-xl">
            <svg
              className="w-6 h-6 text-blue-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
              />
            </svg>
          </div>
          <h2 className="text-xl sm:text-2xl font-bold text-white">
            Desktop Apps
          </h2>
        </div>
        <ResponsiveContainer width="100%" height={350}>
          {renderChart(apps, COLORS.apps, "Apps") || <div />}
        </ResponsiveContainer>
      </div>

      {/* Sites */}
      <div className="group p-6 sm:p-8 bg-gray-800 rounded-3xl shadow-xl ring-1 ring-gray-700 hover:shadow-2xl transition-all duration-300">
        <div className="flex items-center gap-3 mb-6">
          <div className="p-2.5 bg-green-900/30 rounded-xl">
            <svg
              className="w-6 h-6 text-green-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
              />
            </svg>
          </div>
          <h2 className="text-xl sm:text-2xl font-bold text-white">Websites</h2>
        </div>
        <ResponsiveContainer width="100%" height={350}>
          {renderChart(sites, COLORS.sites, "Sites") || <div />}
        </ResponsiveContainer>
      </div>
    </div>
  );
}
