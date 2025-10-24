import type { Metadata } from "next";
import "../globals.css";
import localFont from "next/font/local";

export const metadata: Metadata = {
  title: "Sign In - OSS Archive",
  description: "Sign in to your OSS Archive farmer account",
};

const generalSans = localFont({
  src: [
    {
      path: "../fonts/GeneralSans-Regular.woff2",
      weight: "400",
      style: "normal",
    },
    {
      path: "../fonts/GeneralSans-Medium.woff2",
      weight: "500",
      style: "medium",
    },
    {
      path: "../fonts/GeneralSans-Semibold.woff2",
      weight: "600",
      style: "normal",
    },
    {
      path: "../fonts/GeneralSans-Bold.woff2",
      weight: "800",
      style: "normal",
    },
  ],
  variable: "--font-general-sans",
});

export default function SigninLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div lang="en" className={`${generalSans.variable}`}>
      <div className="overflow-x-hidden font-generalSans leading-[1.25rem] tracking-tight text-black">
        {children}
      </div>
    </div>
  );
}
