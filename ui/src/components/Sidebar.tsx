import Link from 'next/link'
import { useRouter } from 'next/navigation'

const Sidebar = () => {
  const router = useRouter()

  const handleLogout = () => {
    // Clear the auth token
    localStorage.removeItem('authToken')
    // Redirect to login page
    router.push('/user/login')
  }

  return (
    <div className="bg-gray-800 text-white w-64 min-h-screen p-4">
      <nav className="mt-8">
        <ul className="space-y-2">
          <li>
            <Link href="/dashboard/nodes" className="block py-2 px-4 hover:bg-gray-700 rounded">
              Nodes
            </Link>
          </li>
          <li>
            <Link href="/dashboard/metadata" className="block py-2 px-4 hover:bg-gray-700 rounded">
              Metadata
            </Link>
          </li>
          <li>
            <Link href="/dashboard/experiments" className="block py-2 px-4 hover:bg-gray-700 rounded">
              Experiments
            </Link>
          </li>
        </ul>
      </nav>
      <button
        onClick={handleLogout}
        className="mt-8 w-full py-2 px-4 bg-red-600 hover:bg-red-700 rounded"
      >
        Logout
      </button>
    </div>
  )
}

export default Sidebar
