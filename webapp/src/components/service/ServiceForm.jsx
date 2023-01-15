import { Label, TextInput, Button } from 'flowbite-react';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useParams } from 'react-router-dom';

const ServiceForm = () => {
  const { id } = useParams();
  const [service, setService] = useState();

  useEffect(() => {
    if (id) {
      (async () => {
        try {
          const { data } = await axios.get(
            `http://localhost:8000/v1/services/${id}`,
          );
          console.log(data);
          setService(data);
        } catch (e) {
          console.error(e);
        }
      })();
    }
  }, [id]);

  return (
    <div className="flex grow h-fit w-fit">
      <div className="flex flex-col gap-4 min-w-[32rem] p-4 dark:bg-gray-800 rounded">
        <h1 className="text-lg dark:text-white">Edit service</h1>
        <form
          className="flex flex-col gap-4 max-w-xl"
          onSubmit={(e) => e.preventDefault()}
        >
          <div>
            <div className="mb-2 block">
              <Label htmlFor="name" value="Name (optional)" />
            </div>
            <TextInput id="name" name="name" type="text" />
          </div>
          <div>
            <div className="mb-2 block">
              <Label htmlFor="url" value="URL" />
            </div>
            <TextInput id="url" name="url" type="text" />
          </div>
          <div>
            <div className="mb-2 block">
              <Label
                htmlFor="checkInterval"
                value="Check interval in seconds"
              />
            </div>
            <TextInput id="checkInterval" name="checkInterval" type="text" />
          </div>
          <div className="flex flex-row gap-4">
            <Link className="flex grow" to="/dashboard">
              <Button color="gray" className="grow">
                Cancel
              </Button>
            </Link>
            <Button className="grow" type="submit">
              Save
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ServiceForm;
