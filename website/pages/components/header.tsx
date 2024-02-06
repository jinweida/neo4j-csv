import React from 'react';
import Image from "next/image"
import Link from 'next/link';
import { useRouter } from 'next/router';

function Header() {
  const router = useRouter();
  // 根据当前路由设置文本颜色
  const home_color = router.pathname === '/' ? 'text-red-600' : '';
  const targt_color = router.pathname === '/target' ? 'text-red-600' : '';
  const setting_color = router.pathname === '/setting' ? 'text-red-600' : '';
  return (
    <header className="flex justify-between items-center h-20">
        <p className="flex items-start mr-auto text-red-600 text-2xl">
          <Link href="#" className="inline-flex justify-center"><span className="icon-[ion--logo-apple] text-3xl"></span> 用户关系360</Link></p>
        <nav className="flex items-center">
          <Link href="/" className={`inline-flex justify-center ${home_color}`}>
            <span className="icon-[material-symbols--cottage-outline-rounded] text-2xl"></span> 首页</Link>
          <Link href="/target" className={`ml-4 px-4 py-2 rounded flex items-center inline-flex justify-center ${targt_color}`}>
            <span className="icon-[material-symbols--view-sidebar-outline] text-2xl"></span> 目标图谱
          </Link>
          <Link href="#" className={`ml-2 px-4 py-2 rounded flex items-center inline-flex justify-center ${setting_color}`}>
            <span className="icon-[material-symbols--settings-slow-motion-outline] text-2xl"></span> 配置</Link>
        </nav>
    </header>
  );
}

export default Header;