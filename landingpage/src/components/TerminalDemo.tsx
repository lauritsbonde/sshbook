
import React, { useEffect, useState, useRef } from 'react';

interface TerminalDemoProps {
  className?: string;
}

const TerminalDemo: React.FC<TerminalDemoProps> = ({ className = '' }) => {
  const [displayedText, setDisplayedText] = useState('');
  const [cursorIndex, setCursorIndex] = useState(0);
  const [showCursor, setShowCursor] = useState(true);
  const [currentLine, setCurrentLine] = useState(0);
  const contentRef = useRef<HTMLDivElement>(null);
  
  const demoLines = [
    { prompt: '$ ', text: 'sshbook add dev-server', delay: 100 },
    { prompt: '> ', text: 'Name: Development Server', delay: 100 },
    { prompt: '> ', text: 'Group: development', delay: 100 },
    { prompt: '> ', text: 'SSH connection added successfully!', delay: 100 },
    { prompt: '$ ', text: 'sshbook list', delay: 1000 },
    { prompt: '', text: '📖 sshbook connections:', delay: 100 },
    { prompt: '', text: '├─ 📁 development', delay: 100 },
    { prompt: '', text: '│  └─ 🖥️  Development Server (dev-server)', delay: 100 },
    { prompt: '', text: '└─ 📁 production', delay: 100 },
    { prompt: '', text: '   ├─ 🖥️  Web Server (web-01)', delay: 100 },
    { prompt: '', text: '   └─ 🖥️  Database (db-01)', delay: 100 },
    { prompt: '$ ', text: 'sshbook connect dev-server', delay: 1000 },
    { prompt: '', text: '✨ Connecting to dev-server...', delay: 100 },
    { prompt: 'dev-server:~ $ ', text: '_', delay: 100 },
  ];

  useEffect(() => {
    if (currentLine >= demoLines.length) {
      // Reset to beginning after a pause
      const timer = setTimeout(() => {
        setCurrentLine(0);
        setCursorIndex(0);
        setDisplayedText('');
      }, 3000);
      return () => clearTimeout(timer);
    }
    
    const currentLineData = demoLines[currentLine];
    
    if (cursorIndex === 0) {
      setDisplayedText(prev => prev + '\n' + currentLineData.prompt);
      setCursorIndex(currentLineData.prompt.length);
      return;
    }
    
    if (cursorIndex <= currentLineData.prompt.length + currentLineData.text.length) {
      const timer = setTimeout(() => {
        const nextChar = currentLineData.text[cursorIndex - currentLineData.prompt.length];
        setDisplayedText(prev => prev + nextChar);
        setCursorIndex(cursorIndex + 1);
      }, 50); // typing speed
      
      return () => clearTimeout(timer);
    } else {
      const timer = setTimeout(() => {
        setCurrentLine(currentLine + 1);
        setCursorIndex(0);
      }, currentLineData.delay);
      
      return () => clearTimeout(timer);
    }
  }, [currentLine, cursorIndex]);

  // Blinking cursor effect
  useEffect(() => {
    const cursorInterval = setInterval(() => {
      setShowCursor(prev => !prev);
    }, 500);
    
    return () => clearInterval(cursorInterval);
  }, []);

  // Scroll to the bottom when text changes
  useEffect(() => {
    if (contentRef.current) {
      contentRef.current.scrollTop = contentRef.current.scrollHeight;
    }
  }, [displayedText]);
  
  return (
    <div className={`terminal-window ${className}`}>
      <div className="terminal-header">
        <div className="terminal-button bg-red-500"></div>
        <div className="terminal-button bg-yellow-500"></div>
        <div className="terminal-button bg-green-500"></div>
        <span className="ml-2 text-xs text-gray-400">terminal</span>
      </div>
      <div 
        ref={contentRef}
        className="terminal-content font-mono whitespace-pre overflow-y-auto" 
      >
        {displayedText}
        {showCursor && <span className="terminal-cursor"></span>}
      </div>
    </div>
  );
};

export default TerminalDemo;
