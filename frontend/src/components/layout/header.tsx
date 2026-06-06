"use client";

import { useAuth } from "@/contexts/auth-context";

export function Header() {
  const { user } = useAuth();

  return (
    <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6">
      <div />
      <div className="flex items-center gap-3">
        <div className="text-right">
          <p className="text-sm font-medium text-gray-900">{user?.email}</p>
          <p className="text-xs text-gray-500 capitalize">{user?.role}</p>
        </div>
        <div className="w-9 h-9 rounded-full bg-blue-600 flex items-center justify-center text-white text-sm font-medium">
          {user?.email?.charAt(0).toUpperCase()}
        </div>
      </div>
    </header>
  );
}
