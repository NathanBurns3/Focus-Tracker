import {
  Bar,
  BarChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

interface Usage {
  name: string;
  minutes_used: number;
  source: string;
}

export default function UsageChart({ data }: { data: Usage[] }) {
  const apps = data.filter((d) => d.source === "desktop");
  const sites = data.filter((d) => d.source === "chrome");

  const customTooltip = {
    formatter: (value: number) => [`${value}`, "Minutes Used"],
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
      {/* Apps */}
      <div className="p-4 bg-white rounded-2xl shadow">
        <h2 className="text-xl font-semibold mb-4 text-gray-800">Apps</h2>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={apps}>
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip formatter={customTooltip.formatter} />
            <Bar dataKey="minutes_used" fill="#4F46E5" />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Sites */}
      <div className="p-4 bg-white rounded-2xl shadow">
        <h2 className="text-xl font-semibold mb-4 text-gray-800">Sites</h2>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={sites}>
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip formatter={customTooltip.formatter} />
            <Bar dataKey="minutes_used" fill="#10B981" />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
