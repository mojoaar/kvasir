export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-full">
      <aside className="w-64 border-r p-4">
        <nav className="space-y-2">
          <div className="font-semibold text-sm">Kvasir</div>
        </nav>
      </aside>
      <main className="flex-1 p-6">{children}</main>
    </div>
  );
}
