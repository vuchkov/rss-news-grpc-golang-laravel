
1. Set up a new Laravel project:
```
composer create-project --prefer-dist laravel/laravel laravel-rss-reader
cd laravel-rss-reader
```

2. Database setup:

3. Configure your database connection in `.env`.

4. Create Models and Migrations:
```
php artisan make:model RssFeed -m
php artisan make:model RssPost -m
```

- `database/migrations/<timestamp>_create_rss_feeds_table.php`:

```
<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('rss_feeds', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->string('url');
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('rss_feeds');
    }
};
```

- `database/migrations/<timestamp>_create_rss_posts_table.php`:
```php
<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('rss_posts', function (Blueprint $table) {
            $table->id();
            $table->foreignId('rss_feed_id')->constrained()->onDelete('cascade');
            $table->string('title');
            $table->string('link');
            $table->text('description')->nullable();
            $table->timestamp('published_at')->nullable();
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('rss_posts');
    }
};
```

4. Run Migrations:
```
php artisan migrate
```

5. Models (`app/Models/RssFeed.php` and `app/Models/RssPost.php`):
6. 
- `app/Models/RssFeed.php`:
```php
<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

class RssFeed extends Model
{
    use HasFactory;

    protected $fillable = ['name', 'url'];

    public function posts(): HasMany
    {
        return $this->hasMany(RssPost::class);
    }
}
```

- `app/Models/RssPost.php`:
```php
<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class RssPost extends Model
{
    use HasFactory;

    protected $fillable = ['rss_feed_id', 'title', 'link', 'description', 'published_at'];
    protected $casts = [
        'published_at' => 'datetime',
    ];
}
```

6. Create Controllers:
```
php artisan make:controller RssFeedController
php artisan make:controller RssPostController
```

- `app/Http/Controllers/RssFeedController.php`:
```
<?php

namespace App\Http\Controllers;

use App\Models\RssFeed;
use Illuminate\Http\Request;

class RssFeedController extends Controller
{
    public function index()
    {
        $feeds = RssFeed::all();
        return view('rss_feeds.index', compact('feeds'));
    }

    public function create()
    {
        return view('rss_feeds.create');
    }

    public function store(Request $request)
    {
        RssFeed::create($request->validate([
            'name' => 'required',
            'url' => 'required|url',
        ]));

        return redirect()->route('rss_feeds.index');
    }

        public function edit(RssFeed $rssFeed)
    {
        return view('rss_feeds.edit', compact('rssFeed'));
    }

    public function update(Request $request, RssFeed $rssFeed)
    {
        $rssFeed->update($request->validate([
            'name' => 'required',
            'url' => 'required|url',
        ]));

        return redirect()->route('rss_feeds.index');
    }

    public function destroy(RssFeed $rssFeed)
    {
        $rssFeed->delete();
        return redirect()->route('rss_feeds.index');
    }
}
```

- `app/Http/Controllers/RssPostController.php`:
```php
<?php

namespace App\Http\Controllers;

use App\Models\RssPost;
use Illuminate\Http\Request;

class RssPostController extends Controller
{
    public function index(Request $request)
    {
        $posts = RssPost::query();

        if ($request->has('feed')) {
            $posts->where('rss_feed_id', $request->input('feed'));
        }

        if ($request->has('date')) {
            $posts->whereDate('published_at', $request->input('date'));
        }

        if ($request->has('title')) {
            $posts->where('title', 'like', '%' . $request->input('title') . '%');
        }

        $posts = $posts->orderBy('published_at', 'desc')->paginate(10);

        return view('rss_posts.index', compact('posts'));
    }
}
```

7. Create Views (`resources/views/rss_feeds` and `resources/views/rss_posts`):

Create `index.blade.php`, `create.blade.php`, and `edit.blade.php` 
in `resources/views/rss_feeds` and `index.blade.php` in `resources/views/rss_posts`. 
These will contain the HTML for the views. 
I show basic example of the `index.blade.php` for feeds here:

- `resources/views/rss_feeds/index.blade.php`:
```
@extends('layouts.app')

@section('content')
    <h1>RSS Feeds</h1>

    <a href="{{ route('rss_feeds.create') }}" class="btn btn-primary">Add New Feed</a>

    <table class="table">
        <thead>
            <tr>
                <th>Name</th>
                <th>URL</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            @foreach ($feeds as $feed)
                <tr>
                    <td>{{ $feed->name }}</td>
                    <td>{{ $feed->url }}</td>
                    <td>
                        <a href="{{ route('rss_feeds.edit', $feed) }}" class="btn btn-sm btn-warning">Edit</a>
                        <form action="{{ route('rss_feeds.destroy', $feed) }}" method="POST" style="display: inline;">
                            @csrf
                            @method('DELETE')
                            <button type="submit" class="btn btn-sm btn-danger" onclick="return confirm('Are you sure?')">Delete</button>
                        </form>
                    </td>
                </tr>
            @endforeach
        </tbody>
    </table>
@endsection
```

8. Routes (`routes/web.php`):
```php
<?php

use App\Http\Controllers\RssFeedController;
use App\Http\Controllers\RssPostController;
use Illuminate\Support\Facades\Route;

Route::resource('rss_feeds', RssFeedController::class);
Route::get('/posts', [RssPostController::class, 'index'])->name('posts.index');
```

## Implement the background job and connect the Laravel application to the Golang RSS Reader service.

1. Create the Job:
```
php artisan make:job FetchRssFeeds
```

2. `app/Jobs/FetchRssFeeds.php`:
```
<?php

namespace App\Jobs;

use App\Models\RssFeed;
use App\Models\RssPost;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldBeUnique;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;
use Throwable;

class FetchRssFeeds implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /**
     * Create a new job instance.
     */
    public function __construct()
    {
        //
    }

    /**
     * Execute the job.
     */
    public function handle(): void
    {
        $feeds = RssFeed::all();

        foreach ($feeds as $feed) {
            try {
                $response = Http::post('http://localhost:8080/parse', [ // Replace with your service URL
                    'urls' => [$feed->url],
                ]);

                if ($response->successful()) {
                    $data = $response->json();
                    if(isset($data['items'])){
                        foreach ($data['items'] as $item) {

                            RssPost::updateOrCreate(
                                [
                                    'rss_feed_id' => $feed->id,
                                    'link' => $item['link'],
                                ],
                                [
                                    'title' => $item['title'],
                                    'description' => $item['description'],
                                    'published_at' => $item['publish_date'],
                                ]
                            );
                        }
                    } else{
                        Log::error("Invalid response format from RSS service for feed: " . $feed->url . " Response: " . $response->body());
                    }


                } else {
                    Log::error("Failed to fetch RSS feed: " . $feed->url . " Status code: " . $response->status() . " Response: " . $response->body());
                }
            } catch (Throwable $e) {
                Log::error("Error fetching RSS feed: " . $feed->url . " Error: " . $e->getMessage());
            }
        }
    }
}
```

3. Dispatch the Job (e.g., in a command or scheduler):

Using a command:
```
php artisan make:command FetchRss
```
- `app/Console/Commands/FetchRss.php`:
```php
<?php

namespace App\Console\Commands;

use App\Jobs\FetchRssFeeds;
use Illuminate\Console\Command;

class FetchRss extends Command
{
    protected $signature = 'rss:fetch';
    protected $description = 'Fetch RSS feeds';

    public function handle()
    {
        FetchRssFeeds::dispatch();
        $this->info('RSS feed fetching job dispatched.');
        return Command::SUCCESS;
    }
}
```

- Run `php artisan rss:fetch` to dispatch the job.

- Using the scheduler (recommended for regular updates):

In `app/Console/Kernel.php`:
```
protected function schedule(Schedule $schedule)
{
    $schedule->command('rss:fetch')->hourly(); // Run every hour
}
```
- Then run `php artisan schedule:run` (or set up a cron job for this command on your server).

4. Install `guzzlehttp/guzzle`:
```
composer require guzzlehttp/guzzle
```

Key improvements and explanations:
- Error Handling: The try-catch block and logging provide better error handling. Now, if fetching or parsing fails for one feed, it won't stop the process for other feeds. The error messages are more informative.
- HTTP Requests with Http Facade: Using Laravel's Http facade is cleaner and more convenient than using Guzzle directly in most cases.
- updateOrCreate method: This method prevents duplicate posts from being created. It updates the post if it already exists (based on rss_feed_id and link), or creates a new one if it doesn't.
- Checking for the "items" key: The code now verifies that the "items" key exists in the JSON response from the Golang service. This handles cases where the service might return an error in a different format.
- Logging Invalid Response: If the response from the Go service is not in the expected format, the code logs an error message with the full response body. This helps in debugging issues with the service.
- Job dispatching via command: Now you can manually dispatch the job using a command, which is useful for testing and one-time updates.
- Scheduler integration: The scheduler now uses the command, which is the recommended way to schedule jobs in Laravel.

5. Update the Views:

Update the views to display the fetched posts. 

- Here's an example for `resources/views/rss_posts/index.blade.php`:
```
@extends('layouts.app')

@section('content')
    <h1>RSS Posts</h1>

    <form action="{{ route('posts.index') }}" method="GET">
        <select name="feed">
            <option value="">All Feeds</option>
            @foreach (\App\Models\RssFeed::all() as $feed)
                <option value="{{ $feed->id }}" {{ request('feed') == $feed->id ? 'selected' : '' }}>{{ $feed->name }}</option>
            @endforeach
        </select>
        <input type="date" name="date" value="{{ request('date') }}">
        <input type="text" name="title" placeholder="Search by title" value="{{ request('title') }}">
        <button type="submit">Filter</button>
    </form>

    <ul>
        @foreach ($posts as $post)
            <li>
                <a href="{{ $post->link }}" target="_blank">{{ $post->title }}</a>
                <p>{{ $post->published_at }}</p>
            </li>
        @endforeach
    </ul>

    {{ $posts->links() }}
@endsection
```

- This completes the integration between the Laravel application and the Golang service. 
- You now have a fully functional RSS reader application with background updates, filtering, and pagination. 
- Remember to create the missing views (`create.blade.php`, `edit.blade.php`) for the RSS feeds management.


