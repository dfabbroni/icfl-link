'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { FaTachometerAlt, FaDatabase, FaFlask, FaSignOutAlt, FaBars, FaTimes } from 'react-icons/fa'
import { FaCircleNodes, FaNode } from 'react-icons/fa6'
import { useState, useEffect } from 'react'

const Sidebar = () => {
  const router = useRouter()
  const [isOpen, setIsOpen] = useState(true)
  const [isMobile, setIsMobile] = useState(false)

  // Check if we're on mobile when component mounts and when window resizes
  useEffect(() => {
    const checkIfMobile = () => {
      setIsMobile(window.innerWidth < 768)
      // Auto-close sidebar on mobile
      if (window.innerWidth < 768) {
        setIsOpen(false)
      } else {
        setIsOpen(true)
      }
    }

    checkIfMobile()
    
    window.addEventListener('resize', checkIfMobile)
    
    return () => window.removeEventListener('resize', checkIfMobile)
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('authToken')
    router.push('/user/login')
  }

  return (
    <>
      {/* Mobile toggle button - only visible on small screens */}
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className="md:hidden fixed top-4 left-4 z-20 bg-gray-800 text-white p-2 rounded-md shadow-lg"
        aria-label={isOpen ? "Close sidebar" : "Open sidebar"}
      >
        {isOpen ? <FaTimes /> : <FaBars />}
      </button>
      
      {/* Sidebar */}
      <div className={`
        bg-gray-900 text-white min-h-screen p-6 flex flex-col justify-between
        fixed md:static z-10 transition-all duration-300 ease-in-out
        ${isOpen ? 'w-64 left-0' : 'w-0 -left-96 md:w-16 md:left-0'}
        ${isMobile && isOpen ? 'shadow-2xl' : ''}
      `}>
        <div className={`${!isOpen && 'md:invisible'} overflow-hidden`}>
          <h2 className="text-2xl font-semibold mb-6 text-center truncate">ICFL Link</h2>
          <nav className="mt-8">
            <ul className="space-y-4">
              <li>
                <Link href="/dashboard" className="flex items-center py-2 px-4 hover:bg-gray-700 rounded">
                  <FaTachometerAlt className="mr-3 flex-shrink-0" />
                  <span className="truncate">Dashboard</span>
                </Link>
              </li>
              <li>
                <Link href="/dashboard/nodes" className="flex items-center py-2 px-4 hover:bg-gray-700 rounded">
                  <FaCircleNodes className="mr-3 flex-shrink-0" />
                  <span className="truncate">Nodes</span>
                </Link>
              </li>
              <li>
                <Link href="/dashboard/metadata" className="flex items-center py-2 px-4 hover:bg-gray-700 rounded">
                  <FaDatabase className="mr-3 flex-shrink-0" />
                  <span className="truncate">Metadata</span>
                </Link>
              </li>
              <li>
                <Link href="/dashboard/experiments" className="flex items-center py-2 px-4 hover:bg-gray-700 rounded">
                  <FaFlask className="mr-3 flex-shrink-0" />
                  <span className="truncate">Experiments</span>
                </Link>
              </li>
            </ul>
          </nav>
        </div>
        
        {/* Always visible icons for collapsed state on desktop */}
        <div className={`hidden md:${!isOpen ? 'flex' : 'hidden'} flex-col items-center mt-8 space-y-6`}>
          <Link href="/dashboard" className="p-2 hover:bg-gray-700 rounded" title="Dashboard">
            <FaTachometerAlt className="text-xl" />
          </Link>
          <Link href="/dashboard/nodes" className="p-2 hover:bg-gray-700 rounded" title="Nodes">
            <FaNode className="text-xl" />
          </Link>
          <Link href="/dashboard/metadata" className="p-2 hover:bg-gray-700 rounded" title="Metadata">
            <FaDatabase className="text-xl" />
          </Link>
          <Link href="/dashboard/experiments" className="p-2 hover:bg-gray-700 rounded" title="Experiments">
            <FaFlask className="text-xl" />
          </Link>
        </div>
        
        <button
          onClick={handleLogout}
          className={`
            py-2 px-4 bg-red-600 hover:bg-red-700 rounded flex items-center justify-center
            ${!isOpen ? 'md:w-10 md:h-10 md:px-0 md:mx-auto' : 'w-full'}
          `}
          title="Logout"
        >
          <FaSignOutAlt className={isOpen ? 'mr-2' : ''} />
          {isOpen && <span>Logout</span>}
        </button>
      </div>
      
      {/* Overlay to close sidebar when clicking outside on mobile */}
      {isMobile && isOpen && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-50 z-0"
          onClick={() => setIsOpen(false)}
        />
      )}
    </>
  )
}

export default Sidebar