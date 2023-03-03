# Quiz

## Env file setup

Create and edit the `.env` files:

<details>
  <summary>Main (<code>.env</code>)</summary><br>

  ```bash
  cp .env.example .env
  ```

  - Change `PORT` to the port you want the app to be hosted on
  - Change `UID` and `GID` to the user and group you want to run the Docker containers as
</details>

<details>
  <summary>Frontend (<code>frontend/.env</code>)</summary><br>

  ```bash
  cd frontend
  cp .env.example .env
  ```

  - Change `PUBLIC_BACKEND` to point to the public location where you host the backend
    - e.g., if hosting locally, use `http://localhost:${PORT}`
  - If you use a different port publicly than the one in your main `.env` file, change `${PORT}` to your public port
</details>

## Main setup + starting the app

See [backend README](./backend/README.md).
After following those instructions,

```bash
docker compose up
```

to build & start the PocketBase backend and build & serve the SvelteKit frontend

Alternatively, if you do not wish to use Docker, do

```bash
cd frontend && npx pnpm install && npm run build && npm run backend
```

to manually build & serve the site.
