
import React from 'react';

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-950 border-t border-gray-800 py-12">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="mb-6 md:mb-0">
            <div className="flex items-center">
              <span className="text-terminal-highlight font-mono font-bold text-xl">✨ sshbook</span>
            </div>
            <p className="text-gray-400 mt-2 max-w-xs">
              Your SSH connections, organized. Everywhere.
            </p>
          </div>
          
          <div className="grid grid-cols-2 gap-8 md:gap-12">
            <div>
              <h3 className="text-white font-semibold mb-3">Links</h3>
              <ul className="space-y-2">
                <li><a href="#" className="text-gray-400 hover:text-terminal-highlight">Home</a></li>
                <li><a href="#features" className="text-gray-400 hover:text-terminal-highlight">Features</a></li>
                <li><a href="#why" className="text-gray-400 hover:text-terminal-highlight">Why sshbook?</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="text-white font-semibold mb-3">Resources</h3>
              <ul className="space-y-2">
                <li><a href="https://github.com" target="_blank" rel="noopener noreferrer" className="text-gray-400 hover:text-terminal-highlight">GitHub</a></li>
                <li><a href="#" className="text-gray-400 hover:text-terminal-highlight">Documentation</a></li>
                <li><a href="#" className="text-gray-400 hover:text-terminal-highlight">Releases</a></li>
              </ul>
            </div>
          </div>
        </div>
        
        <div className="border-t border-gray-800 mt-12 pt-6 flex flex-col md:flex-row justify-between items-center">
          <p className="text-gray-500 text-sm">© {new Date().getFullYear()} sshbook. All rights reserved.</p>
          <div className="mt-4 md:mt-0 flex space-x-4">
            <a href="#" className="text-gray-400 hover:text-terminal-highlight">Privacy</a>
            <a href="#" className="text-gray-400 hover:text-terminal-highlight">Terms</a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
