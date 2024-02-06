import "@/styles/globals.css";
import type { AppProps } from "next/app";
import Header from "./components/header"
import { useRouter } from 'next/router';

export default function App({ Component, pageProps }: AppProps) {
  const router = useRouter();

  if (router.pathname === '/login') {
    return <Component {...pageProps} />;
  }
  return (
    <div className="container mx-auto">
      <Header />
      <Component {...pageProps} />
    </div>);
}
