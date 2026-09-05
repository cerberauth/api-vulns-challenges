import { ApolloServer } from '@apollo/server';
import { expressMiddleware } from '@apollo/server/express4';
import { ApolloServerPluginDrainHttpServer } from '@apollo/server/plugin/drainHttpServer';
import cors from 'cors';
import express from 'express';
import helmet from 'helmet';
import http from 'http';
import morgan from 'morgan';

// A schema is a collection of type definitions (hence "typeDefs")
// that together define the "shape" of queries that are executed against
// your data.
const typeDefs = `#graphql
  # Comments in GraphQL strings (such as this one) start with the hash (#) symbol.

  # This "Book" type defines the queryable fields for every book in our data source.
  type Book {
    title: String
    author: String
  }

  # The "Query" type is special: it lists all of the available queries that
  # clients can execute, along with the return type for each. In this
  # case, the "books" query returns an array of zero or more Books (defined above).
  type Query {
    books: [Book]
  }
`;

const books = [
  {
    title: 'The Awakening',
    author: 'Kate Chopin',
  },
  {
    title: 'City of Glass',
    author: 'Paul Auster',
  },
];

// Resolvers define how to fetch the types defined in your schema.
// This resolver retrieves books from the "books" array above.
const resolvers = {
  Query: {
    books: () => books,
  },
};

// Toggle between the vulnerable and the fixed, non-vulnerable configuration.
// Defaults to the vulnerable mode, matching the other challenges in this
// repository. Set VULNERABLE=false to run the fixed configuration.
const vulnerable = process.env.VULNERABLE !== 'false';

const app = express();
const httpServer = http.createServer(app);

const logger = morgan(function (tokens, req, res) {
  return [
    tokens.method(req, res),
    tokens.url(req, res),
    tokens.status(req, res),
    tokens.res(req, res, 'content-length'),
    JSON.stringify((req as express.Request).body || ''),
    '-',
    tokens['response-time'](req, res),
    'ms',
  ].join(' ');
});

// The ApolloServer constructor requires two parameters: your schema
// definition and your set of resolvers.
const server = new ApolloServer({
  typeDefs,
  resolvers,
  plugins: [ApolloServerPluginDrainHttpServer({ httpServer })],
  // vulnerable: introspection lets anyone dump the full schema, including
  // fields and types never meant to be discoverable by a client
  introspection: vulnerable,
});

await server.start();

app.use(
  '/graphql',
  // vulnerable: any origin is allowed to make credentialed cross-site
  // requests to the GraphQL endpoint
  vulnerable
    ? cors<cors.CorsRequest>()
    : cors<cors.CorsRequest>({ origin: 'https://trusted.example.com' }),
  helmet(),
  express.json(),
  logger,
  expressMiddleware(server, {
    context: async ({ req }) => ({ token: req.headers.token }),
  }),
);

await new Promise<void>((resolve) =>
  httpServer.listen({ port: 4000 }, resolve),
);
console.log(`🚀 Server ready at http://localhost:4000/`);
