import { Navbar } from "@/components/Navbar"

export function HomePage() {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />
      <main className="flex flex-1 items-center justify-center">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Welcome to ChatApp</h1>
          <p className="mt-2 text-muted-foreground">Your chat application is ready!</p>
        </div>
      </main>
    </div>
  )
}