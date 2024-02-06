/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  distDir: 'build',
  async rewrites() {
    return [
      {
        source: '/sso/login',
        destination: 'http://neo4j:8000/sso/login',
      },
      {
        source: '/api/target_graph',
        destination: 'http://neo4j:8000/api/target_graph',
      },
    ]
  },
};

export default nextConfig;
