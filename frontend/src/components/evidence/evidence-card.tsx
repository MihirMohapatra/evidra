import Link from "next/link";

const statusColors: Record<string, string> = {
  draft: "bg-gray-100 text-gray-700",
  needs_review: "bg-yellow-100 text-yellow-700",
  approved: "bg-green-100 text-green-700",
  rejected: "bg-red-100 text-red-700",
  exported: "bg-blue-100 text-blue-700",
};

export function EvidenceCard({ item }: { item: any }) {
  return (
    <Link
      href={`/evidence/${item.id}`}
      className="block bg-white rounded-xl border border-gray-200 p-5 hover:shadow-md transition-shadow"
    >
      <div className="flex items-start justify-between mb-3">
        <h3 className="font-medium text-gray-900">{item.title}</h3>
        <span
          className={`text-xs font-medium px-2 py-1 rounded-full ${
            statusColors[item.status] || "bg-gray-100 text-gray-700"
          }`}
        >
          {item.status?.replace("_", " ")}
        </span>
      </div>
      <p className="text-sm text-gray-500 line-clamp-2 mb-3">{item.content}</p>
      <div className="flex items-center gap-2 text-xs text-gray-400">
        <span className="bg-gray-100 px-2 py-0.5 rounded">{item.category}</span>
        {item.tags?.slice(0, 3).map((tag: string) => (
          <span key={tag} className="bg-gray-100 px-2 py-0.5 rounded">
            {tag}
          </span>
        ))}
      </div>
    </Link>
  );
}
